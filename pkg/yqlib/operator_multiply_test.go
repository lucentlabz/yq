package yqlib

import (
	"testing"
)

// crossfunction:
// normal single document, treat each match separately: eval
// multiple documents, how do you do things between docs??? eval-all
// alternative
// collect and filter matches according to doc before running expression
// d0a d1a d0b d1b
// run it for (d0a d1a) x (d0b d1b) - noting that we dont do (d0a x d1a) nor (d0b d1b)
// I think this will work for eval-all correctly then..

//alternative, like jq, eval-all puts all docs in an single array.
var multiplyOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		document:   `{a: {also: [1]}, b: {also: me}}`,
		expression: `. * {"a" : .b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: me}, b: {also: me}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: {also: [1]}`,
		document2:  `b: {also: me}`,
		expression: `.a * .b`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: me}, b: {also: me}}\n",
		},
	},
	{
		skipDoc:    true,
		expression: `{} * {"cat":"dog"}`,
		expected: []string{
			"D0, P[], (!!map)::cat: dog\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {also: me}, b: {also: [1]}}`,
		expression: `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: [1]}, b: {also: [1]}}\n",
		},
	},
	{
		description: "Merge objects together, returning merged result only",
		document:    `{a: {field: me, fieldA: cat}, b: {field: {g: wizz}, fieldB: dog}}`,
		expression:  `.a * .b`,
		expected: []string{
			"D0, P[a], (!!map)::{field: {g: wizz}, fieldA: cat, fieldB: dog}\n",
		},
	},
	{
		description: "Merge objects together, returning parent object",
		document:    `{a: {field: me, fieldA: cat}, b: {field: {g: wizz}, fieldB: dog}}`,
		expression:  `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {field: {g: wizz}, fieldA: cat, fieldB: dog}, b: {field: {g: wizz}, fieldB: dog}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {also: {g: wizz}}, b: {also: me}}`,
		expression: `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: me}, b: {also: me}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {also: {g: wizz}}, b: {also: [1]}}`,
		expression: `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: [1]}, b: {also: [1]}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {also: [1]}, b: {also: {g: wizz}}}`,
		expression: `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {also: {g: wizz}}, b: {also: {g: wizz}}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {things: great}, b: {also: me}}`,
		expression: `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {things: great, also: me}, b: {also: me}}\n",
		},
	},
	{
		description:           "Merge keeps style of LHS",
		dontFormatInputForDoc: true,
		document: `a: {things: great}
b:
  also: "me"
`,
		expression: `. * {"a":.b}`,
		expected: []string{
			`D0, P[], (!!map)::a: {things: great, also: "me"}
b:
    also: "me"
`,
		},
	},
	{
		description: "Merge arrays",
		document:    `{a: [1,2,3], b: [3,4,5]}`,
		expression:  `. * {"a":.b}`,
		expected: []string{
			"D0, P[], (!!map)::{a: [3, 4, 5], b: [3, 4, 5]}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [1], b: [2]}`,
		expression: `.a *+ .b`,
		expected: []string{
			"D0, P[a], (!!seq)::[1, 2]\n",
		},
	},
	{
		description: "Merge, only existing fields",
		document:    `{a: {thing: one, cat: frog}, b: {missing: two, thing: two}}`,
		expression:  `.a *? .b`,
		expected: []string{
			"D0, P[a], (!!map)::{thing: two, cat: frog}\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: [{thing: one}], b: [{missing: two, thing: two}]}`,
		expression: `.a *? .b`,
		expected: []string{
			"D0, P[a], (!!seq)::[{thing: two}]\n",
		},
	},
	{
		description: "Merge, appending arrays",
		document:    `{a: {array: [1, 2, animal: dog], value: coconut}, b: {array: [3, 4, animal: cat], value: banana}}`,
		expression:  `.a *+ .b`,
		expected: []string{
			"D0, P[a], (!!map)::{array: [1, 2, {animal: dog}, 3, 4, {animal: cat}], value: banana}\n",
		},
	},
	{
		description: "Merge, only existing fields, appending arrays",
		document:    `{a: {thing: [1,2]}, b: {thing: [3,4], another: [1]}}`,
		expression:  `.a *?+ .b`,
		expected: []string{
			"D0, P[a], (!!map)::{thing: [1, 2, 3, 4]}\n",
		},
	},
	{
		description: "Merge to prefix an element",
		document:    `{a: cat, b: dog}`,
		expression:  `. * {"a": {"c": .a}}`,
		expected: []string{
			"D0, P[], (!!map)::{a: {c: cat}, b: dog}\n",
		},
	},
	{
		description: "Merge with simple aliases",
		document:    `{a: &cat {c: frog}, b: {f: *cat}, c: {g: thongs}}`,
		expression:  `.c * .b`,
		expected: []string{
			"D0, P[c], (!!map)::{g: thongs, f: *cat}\n",
		},
	},
	{
		description: "Merge copies anchor names",
		document:    `{a: {c: &cat frog}, b: {f: *cat}, c: {g: thongs}}`,
		expression:  `.c * .a`,
		expected: []string{
			"D0, P[c], (!!map)::{g: thongs, c: &cat frog}\n",
		},
	},
	{
		description: "Merge with merge anchors",
		document:    mergeDocSample,
		expression:  `.foobar * .foobarList`,
		expected: []string{
			"D0, P[foobar], (!!map)::c: foobarList_c\n<<: [*foo, *bar]\nthing: foobar_thing\nb: foobarList_b\n",
		},
	},
}

func TestMultiplyOperatorScenarios(t *testing.T) {
	for _, tt := range multiplyOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Multiply", multiplyOperatorScenarios)
}
