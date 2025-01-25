package regexp_P

import (
	"regexp"

	gleam_P "example.com/todo/gleam"
	option_P "example.com/todo/gleam_stdlib/gleam/option"
)

type Regexp_t struct {
	re *regexp.Regexp
}

func (r Regexp_t) Hash() uint32 {
	return gleam_P.String_t(r.re.String()).Hash()
}

func (r Regexp_t) Equal(o Regexp_t) bool {
	return r.re.String() == o.re.String()
}

func doCheck(re Regexp_t, s gleam_P.String_t) gleam_P.Bool_t {
	return gleam_P.Bool_t(re.re.MatchString(string(s)))
}

func doCompile(pattern gleam_P.String_t, with Options_t) gleam_P.Result_t[Regexp_t, CompileError_t] {
	re, err := regexp.Compile(string(pattern))
	if err != nil {
		return gleam_P.Error_c[Regexp_t, CompileError_t]{CompileError_c{gleam_P.String_t(err.Error()), -1}}
	}
	flags := ""
	if with.CaseInsensitive {
		flags += "i"
	}
	if with.MultiLine {
		flags += "m"
	}
	if flags != "" {
		pattern = gleam_P.String_t("(?" + flags + ")" + string(pattern))
	}
	return gleam_P.Ok_c[Regexp_t, CompileError_t]{Regexp_t{re}}
}

func doSplit(re Regexp_t, s gleam_P.String_t) gleam_P.List_t[gleam_P.String_t] {
	parts := re.re.Split(string(s), -1)
	elems := make([]gleam_P.String_t, len(parts))
	for i, part := range parts {
		elems[i] = gleam_P.String_t(part)
	}
	return gleam_P.ToList(elems...)
}

func doScan(re Regexp_t, s gleam_P.String_t) gleam_P.List_t[Match_t] {
	matches := re.re.FindAllStringSubmatch(string(s), -1)
	res := make([]Match_t, len(matches))
	for i, match := range matches {
		content := match[0]
		submatches := make([]option_P.Option_t[gleam_P.String_t], len(match)-1)
		for n := len(match) - 1; n > 0; n-- {
			if match[n] != "" {
				submatches[n-1] = option_P.Some_c[gleam_P.String_t]{gleam_P.String_t(match[n])}
				continue
			}
			if len(submatches) > 0 {
				submatches[n-1] = option_P.None_c[gleam_P.String_t]{}
			}
		}
		res[i] = Match_c{gleam_P.String_t(content), gleam_P.ToList(submatches...)}
	}
	return gleam_P.ToList(res...)
}

func Replace(each Regexp_t, in gleam_P.String_t, with gleam_P.String_t) gleam_P.String_t {
	return gleam_P.String_t(each.re.ReplaceAllString(string(in), string(with)))
}
