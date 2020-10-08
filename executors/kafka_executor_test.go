package executors

import (
	"clamp-core/config"
	"testing"
)

func TestKafkaVal_DoExecute(t *testing.T) {
	type fields struct {
		ConnectionURL string
		TopicName     string
		ContentType   string
		ReplyTo       string
	}
	type args struct {
		requestBody interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "should return error if connection fails",
			fields: fields{
<<<<<<< HEAD
				ConnectionURL: "localhost:19092/",
=======
				ConnectionURL: "54.190.25.178:19092/",
>>>>>>> master
				TopicName:     "topic_test",
				ContentType:   "text/plain",
			},
			args: args{
				requestBody: "message body",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should publish message to kafka topic",
			fields: fields{
				ConnectionURL: config.ENV.KafkaConnectionStr,
				TopicName:     config.ENV.KafkaTopicName,
				ContentType:   "text/plain",
			},
			args: args{
				requestBody: "message body",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := KafkaVal{
				ConnectionURL: tt.fields.ConnectionURL,
				TopicName:     tt.fields.TopicName,
				ContentType:   tt.fields.ContentType,
			}
			got, err := val.DoExecute(tt.args.requestBody, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("DoExecute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				t.Errorf("DoExecute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
