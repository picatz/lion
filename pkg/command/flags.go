package command

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Flag struct {
	Name     string
	Desc     string
	Optional bool
	Default  interface{}
	Value    interface{}
}

type Flags map[string]*Flag

func (c Flags) HelpString() string {
	var (
		lineTemplate string
		longestFlag  int = 6 // --help
		allFlags         = map[string]string{}
		all              = []string{}
	)

	for _, v := range c {
		var (
			flagStr  = fmt.Sprintf("--%s", v.Name)
			flagDesc = v.Desc
		)
		allFlags[flagStr] = flagDesc
		if len(flagStr) > longestFlag {
			longestFlag = len(flagStr)
		}
	}

	longestFlag += 2 // include just a little extra padding

	lineTemplate = "\t%-" + fmt.Sprintf("%d", longestFlag) + "s" + " %s \n"

	for flag, desc := range allFlags {
		all = append(all, fmt.Sprintf(lineTemplate, flag, desc))
	}

	// add help flag to the end
	all = append(all, fmt.Sprintf(lineTemplate, "--help", "Print this help menu (default: false)"))

	return strings.TrimRight(strings.Join(all, ""), "\n")
}

func (c Flags) AllBut(ignore ...string) Flags {
	allBut := Flags{}

	for k, flag := range c {
		if !contains(ignore, k) {
			allBut[k] = flag
		}
	}

	return allBut
}

func (c Flags) GetValue(flag string) (interface{}, error) {
	cmdFlag, ok := c[flag]
	if !ok {
		return nil, fmt.Errorf("the %q flag was not found", flag)
	}
	if cmdFlag.Value == nil {
		return cmdFlag.Default, nil
	}
	return cmdFlag.Value, nil
}

func (c Flags) GetReflectValue(flag string) (reflect.Value, error) {
	i, err := c.GetValue(flag)
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	return reflect.ValueOf(i), nil
}

func (c Flags) StringValue(flag string) (string, error) {
	cmdFlag, ok := c[flag]
	if !ok {
		return "", fmt.Errorf("the %q flag was not found", flag)
	}
	ptr, ok := cmdFlag.Value.(*string)
	if !ok {
		return "", fmt.Errorf("the %q flag is a %T not a *string", flag, cmdFlag.Value)
	}
	if ptr == nil {
		return cmdFlag.Default.(string), nil
	}
	return *ptr, nil
}

func (c Flags) BoolValue(flag string) (bool, error) {
	cmdFlag, ok := c[flag]
	if !ok {
		return false, fmt.Errorf("the %q flag was not found", flag)
	}
	ptr, ok := cmdFlag.Value.(*bool)
	if !ok {
		return false, fmt.Errorf("the %q flag is a %T not *bool", flag, cmdFlag.Value)
	}
	if ptr == nil {
		return cmdFlag.Default.(bool), nil
	}
	return *ptr, nil
}

func (c Flags) IntValue(flag string) (int, error) {
	cmdFlag, ok := c[flag]
	if !ok {
		return 0, fmt.Errorf("the %q flag was not found", flag)
	}
	ptr, ok := cmdFlag.Value.(*int)
	if !ok {
		return 0, fmt.Errorf("the %q flag is a %T not *int", flag, cmdFlag.Value)
	}
	if ptr == nil {
		return cmdFlag.Default.(int), nil
	}
	return *ptr, nil
}

func (c Flags) Int64Value(flag string) (int64, error) {
	cmdFlag, ok := c[flag]
	if !ok {
		return 0, fmt.Errorf("the %q flag was not found", flag)
	}
	ptr, ok := cmdFlag.Value.(*int64)
	if !ok {
		return 0, fmt.Errorf("the %q flag is a %T not *int64", flag, cmdFlag.Value)
	}
	if ptr == nil {
		return cmdFlag.Default.(int64), nil
	}
	return *ptr, nil
}

func (c Flags) DurationValue(flag string) (time.Duration, error) {
	cmdFlag, ok := c[flag]
	if !ok {
		return 0, fmt.Errorf("the %q flag was not found", flag)
	}
	ptr, ok := cmdFlag.Value.(*time.Duration)
	if !ok {
		return 0, fmt.Errorf("the %q flag is a %T not *int64", flag, cmdFlag.Value)
	}
	if ptr == nil {
		return cmdFlag.Default.(time.Duration), nil
	}
	return *ptr, nil
}

func (c Flags) StringSliceValue(flag string) ([]string, error) {
	str, err := c.StringValue(flag)
	if err != nil || strings.TrimSpace(str) == "" {
		return []string{}, err
	}
	parts := strings.Split(str, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts, nil
}

func (c Flags) RegexValue(flag string) (*regexp.Regexp, error) {
	matchStr, err := c.StringValue(flag)
	if err != nil {
		return nil, err
	}
	match, err := regexp.Compile(matchStr)
	if err != nil {
		return nil, err
	}
	return match, nil
}

func (c Flags) TimeValue(flag string) (time.Time, error) {
	str, err := c.StringValue(flag)
	if err != nil {
		return time.Time{}, err
	}
	timeValue, err := time.Parse("Jan-02-2006", str)
	if err != nil {
		return time.Time{}, err
	}
	return timeValue, nil
}
