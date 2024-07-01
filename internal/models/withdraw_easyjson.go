// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson4f4a6fc6DecodeGithubComFlutterDizasterGophermartBonusInternalModels(in *jlexer.Lexer, out *Withdrawals) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Withdrawals, 0, 1)
			} else {
				*out = Withdrawals{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Withdraw
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson4f4a6fc6EncodeGithubComFlutterDizasterGophermartBonusInternalModels(out *jwriter.Writer, in Withdrawals) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Withdrawals) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4f4a6fc6EncodeGithubComFlutterDizasterGophermartBonusInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Withdrawals) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4f4a6fc6EncodeGithubComFlutterDizasterGophermartBonusInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Withdrawals) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4f4a6fc6DecodeGithubComFlutterDizasterGophermartBonusInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Withdrawals) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4f4a6fc6DecodeGithubComFlutterDizasterGophermartBonusInternalModels(l, v)
}
func easyjson4f4a6fc6DecodeGithubComFlutterDizasterGophermartBonusInternalModels1(in *jlexer.Lexer, out *Withdraw) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "order":
			out.StringOrderID = string(in.String())
		case "sum":
			out.Sum = float64(in.Float64())
		case "processed_at":
			out.ProcessedAt = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson4f4a6fc6EncodeGithubComFlutterDizasterGophermartBonusInternalModels1(out *jwriter.Writer, in Withdraw) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"order\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.StringOrderID))
	}
	{
		const prefix string = ",\"sum\":"
		out.RawString(prefix)
		out.Float64(float64(in.Sum))
	}
	if in.ProcessedAt != "" {
		const prefix string = ",\"processed_at\":"
		out.RawString(prefix)
		out.String(string(in.ProcessedAt))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Withdraw) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4f4a6fc6EncodeGithubComFlutterDizasterGophermartBonusInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Withdraw) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4f4a6fc6EncodeGithubComFlutterDizasterGophermartBonusInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Withdraw) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4f4a6fc6DecodeGithubComFlutterDizasterGophermartBonusInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Withdraw) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4f4a6fc6DecodeGithubComFlutterDizasterGophermartBonusInternalModels1(l, v)
}
