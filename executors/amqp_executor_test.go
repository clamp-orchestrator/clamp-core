package executors

import (
	"testing"
)

func TestAMQPVal_DoExecute(t *testing.T) {
	type fields struct {
		ConnectionURL string
		QueueName     string
		ExchangeName  string
		RoutingKey    string
		ContentType   string
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
		/*{
			name: "should send message to queue",
			fields: fields{
				ConnectionURL: "amqp:///",
				QueueName:     "test_queue",
				ExchangeName:  "",
				RoutingKey:    "",
				ContentType:   "text/plain",
			},
			args: args{
				requestBody: "message body",
			},
			want:    nil,
			wantErr: false,
		}, {
			name: "should send message to exchange",
			fields: fields{
				ConnectionURL: "amqp:///",
				QueueName:     "",
				ExchangeName:  "test_exchange",
				RoutingKey:    "test_key",
				ContentType:   "text/plain",
			},
			args: args{
				requestBody: "message body",
			},
			want:    nil,
			wantErr: false,
		},*/{
			name: "should return error if connection fails",
			fields: fields{
				ConnectionURL: "amqp://test:test@localhost:5672/",
				QueueName:     "",
				ExchangeName:  "test_exchange",
				RoutingKey:    "test_key",
				ContentType:   "text/plain",
			},
			args: args{
				requestBody: "message body",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := AMQPVal{
				ConnectionURL: tt.fields.ConnectionURL,
				QueueName:     tt.fields.QueueName,
				ExchangeName:  tt.fields.ExchangeName,
				RoutingKey:    tt.fields.RoutingKey,
				ContentType:   tt.fields.ContentType,
			}
			got, err := val.DoExecute(tt.args.requestBody)
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
