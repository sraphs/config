package config

import (
	"testing"

	"github.com/sraphs/encoding"
	"github.com/stretchr/testify/assert"
)

func TestDescriptor_GetCodec(t *testing.T) {
	type fields struct {
		Name   string
		Format string
		Data   []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   encoding.Codec
	}{
		{
			name: "json",
			fields: fields{
				Name:   "test",
				Format: "json",
				Data:   []byte("{}"),
			},
			want: encoding.GetCodec("json"),
		},
		{
			name: "yaml",
			fields: fields{
				Name:   "test",
				Format: "yaml",
				Data:   []byte(""),
			},
			want: encoding.GetCodec("yaml"),
		},
		{
			name: "yml",
			fields: fields{
				Name:   "test",
				Format: "yml",
				Data:   []byte(""),
			},
			want: encoding.GetCodec("yaml"),
		},
		{
			name: "env",
			fields: fields{
				Name:   "test",
				Format: "env",
				Data:   []byte(""),
			},
			want: encoding.GetCodec("env"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Descriptor{
				Name:   tt.fields.Name,
				Format: tt.fields.Format,
				Data:   tt.fields.Data,
			}
			got := d.GetCodec()
			assert.Equal(t, tt.want, got)
		})
	}
}
