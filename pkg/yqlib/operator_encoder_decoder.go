package yqlib

import (
	"bufio"
	"bytes"
	"container/list"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func configureEncoder(format PrinterOutputFormat, indent int) Encoder {
	switch format {
	case JSONOutputFormat:
		return NewJSONEncoder(indent, false, false)
	case PropsOutputFormat:
		return NewPropertiesEncoder(true)
	case CSVOutputFormat:
		return NewCsvEncoder(',')
	case TSVOutputFormat:
		return NewCsvEncoder('\t')
	case YamlOutputFormat:
		return NewYamlEncoder(indent, false, ConfiguredYamlPreferences)
	case XMLOutputFormat:
		return NewXMLEncoder(indent, ConfiguredXMLPreferences)
	case Base64OutputFormat:
		return NewBase64Encoder()
	}
	panic("invalid encoder")
}

func encodeToString(candidate *CandidateNode, prefs encoderPreferences) (string, error) {
	var output bytes.Buffer

	encoder := configureEncoder(prefs.format, prefs.indent)

	printer := NewPrinter(encoder, NewSinglePrinterWriter(bufio.NewWriter(&output)))
	err := printer.PrintResults(candidate.AsList())
	return output.String(), err
}

type encoderPreferences struct {
	format PrinterOutputFormat
	indent int
}

/* encodes object as yaml string */

func encodeOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	preferences := expressionNode.Operation.Preferences.(encoderPreferences)
	var results = list.New()

	hasOnlyOneNewLine := regexp.MustCompile("[^\n].*\n$")
	endWithNewLine := regexp.MustCompile(".*\n$")
	chomper := regexp.MustCompile("\n+$")

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		stringValue, err := encodeToString(candidate, preferences)

		if err != nil {
			return Context{}, err
		}

		// remove trailing newlines if needed.
		// check if we originally decoded this path, and the original thing had a single line.
		originalList := context.GetVariable("decoded: " + candidate.GetKey())
		if originalList != nil && originalList.Len() > 0 && hasOnlyOneNewLine.MatchString(stringValue) {

			original := originalList.Front().Value.(*CandidateNode)
			originalNode := unwrapDoc(original.Node)
			// original block did not have a new line at the end, get rid of this one too
			if !endWithNewLine.MatchString(originalNode.Value) {
				stringValue = chomper.ReplaceAllString(stringValue, "")
			}
		}

		// dont print a new line when printing json on a single line.
		if (preferences.format == JSONOutputFormat && preferences.indent == 0) ||
			preferences.format == CSVOutputFormat ||
			preferences.format == TSVOutputFormat {
			stringValue = chomper.ReplaceAllString(stringValue, "")
		}

		stringContentNode := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: stringValue}
		results.PushBack(candidate.CreateReplacement(stringContentNode))
	}
	return context.ChildContext(results), nil
}

type decoderPreferences struct {
	format InputFormat
}

/* takes a string and decodes it back into an object */
func decodeOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {

	preferences := expressionNode.Operation.Preferences.(decoderPreferences)

	var decoder Decoder
	switch preferences.format {
	case YamlInputFormat:
		decoder = NewYamlDecoder(ConfiguredYamlPreferences)
	case XMLInputFormat:
		decoder = NewXMLDecoder(ConfiguredXMLPreferences)
	case Base64InputFormat:
		decoder = NewBase64Decoder()
	case PropertiesInputFormat:
		decoder = NewPropertiesDecoder()
	case CSVObjectInputFormat:
		decoder = NewCSVObjectDecoder(',')
	case TSVObjectInputFormat:
		decoder = NewCSVObjectDecoder('\t')
	}

	var results = list.New()
	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		context.SetVariable("decoded: "+candidate.GetKey(), candidate.AsList())


		err := decoder.Init(strings.NewReader(unwrapDoc(candidate.Node).Value))
		if err != nil {
			return Context{}, err
		}

		decodedNode, errorReading := decoder.Decode()
		if errorReading != nil {
			return Context{}, errorReading
		}
		//first node is a doc
		node := unwrapDoc(decodedNode.Node)

		results.PushBack(candidate.CreateReplacement(node))
	}
	return context.ChildContext(results), nil
}
