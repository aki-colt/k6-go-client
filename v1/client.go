package k6cli

type K6Client struct {
	Options K6Options
	Block   K6Block
	Imports []string
}

func (c K6Client) GenScript() (string, error) {
	imports := make(map[string]any)
	block, err := c.Block.genBlockScript(imports)
	if err != nil {
		return "", err
	}
	for k := range imports {
		c.Imports = append(c.Imports, k)
	}
	res := ""
	for _, imp := range c.Imports {
		res += imp + "\n"
	}
	res += "\n"
	res += c.Options.genOptionScript() + "\n"
	res += "\n"
	res += "export default function() {\n" + block + "}"
	return res, nil
}
