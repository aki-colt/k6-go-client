package k6cli

import (
	"encoding/json"
	"fmt"
)

type K6Block interface {
	// genBlockScript generate the js script for the block, imports is a set contained what package is imported by all blocks
	genBlockScript(imports map[string]any) string
}

// K6RawCodeBlock is a block only with raw js code
type K6RawCodeBlock struct {
	Code string
}

func (b K6RawCodeBlock) genBlockScript(imports map[string]any) string {
	return b.Code
}

// normal code block contains constants and variable declares, request declares and inner blocks
// the order that generated script will be is:
// 1. constants declares
// 2. variable declares
// 3. request declares
// 4. inner blocks
type K6NormalBlock struct {
	// constants is the declare of const variables in js in this block
	// Example:
	// map[string]any{"a":20, "b":"111", "c": {"d":1}}
	// will be generated as
	// const a = 20; const b = "111"; const c = {d: 1}
	Constants map[string]any
	// variables is the declare of variables in js in this block
	// Example:
	// map[string]any{"a":20, "b":"111", "c": {"d":1}}
	// will be generated as
	// let a = 20; let b = "111"; let c = {d: 1}
	Variables map[string]any
	Requests  []K6Request
	// blocks is the blocks in the outer block, it can be an if block, a for block, a normal block or a raw code block
	// block will not be added into js script with {}, so it can be seen as code piece
	Blocks []K6Block
}

func (b K6NormalBlock) genBlockScript(imports map[string]any) string {
	res := ""
	if len(b.Constants) != 0 {
		for k, v := range b.Constants {
			vStr, _ := json.Marshal(&v)
			res += fmt.Sprintf("const %s = %s\n", k, vStr)
		}
		res += "\n"
	}
	if len(b.Variables) != 0 {
		for k, v := range b.Variables {
			vStr, _ := json.Marshal(&v)
			res += fmt.Sprintf("let %s = %s\n", k, vStr)
		}
		res += "\n"
	}
	if len(b.Requests) != 0 {
		for _, req := range b.Requests {
			res += req.genRequestScript(imports) + "\n"
		}
		res += "\n"
	}
	if len(b.Blocks) != 0 {
		for _, block := range b.Blocks {
			res += block.genBlockScript(imports) + "\n\n"
		}
	}
	return res
}

type K6IfBlock struct {
	// condition after if in js grammar
	// Example:
	// a === 10
	Condition string
	// the block that will be added into if block
	IfBlock K6NormalBlock
	// the block that will be added into else block
	ElseBlock K6NormalBlock
}

func (b K6IfBlock) genBlockScript(imports map[string]any) string {
	return fmt.Sprintf("if (%s) {\n%s\n} else {\n%s\n}",
		b.Condition, b.IfBlock.genBlockScript(imports), b.ElseBlock.genBlockScript(imports))
}

type K6ForBlock struct {
	// condition after for in js grammar
	// Example:
	// let i = 0; i < 100; i++
	Condition string
	Block     K6NormalBlock
}

func (b K6ForBlock) genBlockScript(imports map[string]any) string {
	return fmt.Sprintf("for (%s) {\n%s\n}", b.Condition, b.Block.genBlockScript(imports))
}
