package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/safejob/dify-sdk-go"
	"github.com/safejob/dify-sdk-go/types"
)

func TestChatflowApp(t *testing.T) {
	client, err := dify.NewClient(dify.ClientConfig{
		ApiServer: os.Getenv("APIServer"),
		ApiKey:    os.Getenv("ChatflowApiKey"),
		User:      "chatflow-demo",
		Debug:     true,
		Timeout:   time.Second * 180,
		Transport: nil,
	})
	if err != nil {
		t.Fatal("初始化客户端失败：", err)
	}

	t.Run("Chatflow_RunBlock", func(t *testing.T) {
		ctx := context.Background()

		input := make(map[string]interface{})
		input["name"] = "张三"

		resp, err := client.ChatflowApp().RunBlock(ctx, types.ChatRequest{
			Query:            "你好!你知道我是谁么？",
			Inputs:           input,
			ResponseMode:     "",
			User:             "",
			ConversationId:   "",
			Files:            nil,
			AutoGenerateName: nil,
		})
		t.Logf("resp=%+v err=%v", resp, err)
	})

	t.Run("Chatflow_Run_ParseToStructCh", func(t *testing.T) {
		ctx := context.Background()

		input := make(map[string]interface{})
		input["name"] = "张三"

		eventCh := client.ChatflowApp().Run(ctx, types.ChatRequest{
			Query:            "帮我构思一个国庆五天的出游计划，尽可能详细一点",
			Inputs:           input,
			ResponseMode:     "",
			User:             "",
			ConversationId:   "",
			Files:            nil,
			AutoGenerateName: nil,
		}).ParseToStructCh()
		for {
			select {
			case msg, ok := <-eventCh:
				if !ok {
					return
				}
				if msg.Event == "error" {
					t.Logf("status=%d code=%s message=%s", msg.Status, msg.Code, msg.Message)
				}
				t.Log(msg.Answer)
			}
		}
	})

	t.Run("Chatflow_Run_SimplePrint", func(t *testing.T) {
		ctx := context.Background()

		input := make(map[string]interface{})
		input["name"] = "张三"

		eventCh := client.DebugOff().ChatflowApp().Run(ctx, types.ChatRequest{
			Query:            "你知道现在的时间以及星期么？",
			Inputs:           input,
			ResponseMode:     "",
			User:             "",
			ConversationId:   "",
			Files:            nil,
			AutoGenerateName: nil,
		}).SimplePrint()
		for {
			select {
			case msg, ok := <-eventCh:
				if !ok {
					return
				}
				fmt.Printf("%s", msg)
			}
		}
	})

	t.Run("Chatflow_Run_ParseToEventCh", func(t *testing.T) {
		ctx := context.Background()

		input := make(map[string]interface{})
		input["name"] = "张三"

		eventCh := client.ChatflowApp().Run(ctx, types.ChatRequest{
			Query:            "你知道现在的时间以及星期么？",
			Inputs:           input,
			ResponseMode:     "",
			User:             "",
			ConversationId:   "",
			Files:            nil,
			AutoGenerateName: nil,
		}).ParseToEventCh()
		for {
			select {
			case msg, ok := <-eventCh:
				if !ok {
					return
				}
				t.Logf("====>event: %s %+v\n", msg.Type, msg.Data)
			}
		}
	})

	t.Run("Chatflow_Run_Stop", func(t *testing.T) {
		ctx := context.Background()

		input := make(map[string]interface{})
		input["name"] = "张三"

		eventCh := client.ChatflowApp().Run(ctx, types.ChatRequest{
			Query:            "帮我构思一个国庆五天的出游计划，尽可能详细一点",
			Inputs:           input,
			ResponseMode:     "",
			User:             "",
			ConversationId:   "",
			Files:            nil,
			AutoGenerateName: nil,
		}).ParseToStructCh()
		cnt := 0
		for {
			select {
			case msg, ok := <-eventCh:
				if !ok {
					return
				}
				t.Log(msg.Answer)
				cnt++
				if cnt == 4 {
					err := client.ChatflowApp().Stop(msg.TaskId, "")
					t.Logf("err=%v", err)
				}
			}
		}
	})

	t.Run("Chatflow_UploadFile", func(t *testing.T) {
		f, err := os.Open("/Users/alsc/Downloads/abcd")
		if err != nil {
			t.Logf("Error opening file: %v\n", err)
			return
		}
		defer f.Close()

		resp, err := client.ChatflowApp().UploadFile(
			"/Users/alsc/Downloads/会议室分布.png",
			nil,
			"",
		)
		t.Logf("resp=%v err=%v\n", resp, err)

		resp, err = client.ChatflowApp().UploadFile(
			"",
			f,
			"",
		)
		t.Logf("resp=%v err=%v\n", resp, err)
	})

	t.Run("Chatflow_AppInfo", func(t *testing.T) {
		resp, err := client.ChatflowApp().AppInfo()
		t.Logf("resp=%+v err=%v", resp, err)
	})

	t.Run("Chatflow_AppParameter", func(t *testing.T) {
		resp, err := client.ChatflowApp().AppParameter()
		t.Logf("resp=%+v err=%v", resp, err)
	})

	t.Run("Chatflow_MsgFeedback", func(t *testing.T) {
		err := client.ChatflowApp().MsgFeedback(types.FeedbackReq{
			MessageId: "a89094dd-8dac-4b51-aa77-920099ae4ef9",
			Rating:    types.MsgFeedbackNull,
			User:      "",
			Content:   "非常不错",
		})
		t.Logf("err=%v", err)
	})

	t.Run("Chatflow_SuggestQuestionList", func(t *testing.T) {
		resp, err := client.ChatflowApp().SuggestQuestionList("c71918e4-bb23-4ff9-bb63-e5fa5aaf6afa", "")
		t.Logf("resp=%v err=%v", resp, err)
	})

	t.Run("Chatflow_History", func(t *testing.T) {
		resp, err := client.ChatflowApp().HistoryPro("0a9a0917-0c36-4121-8934-17367bb803c0", "", "", 20)
		t.Logf("resp=%v err=%v", resp, err)
	})

	t.Run("Chatflow_ConversationList", func(t *testing.T) {
		resp, err := client.ChatflowApp().ConversationList("")
		t.Logf("resp=%v err=%v", resp, err)
	})

	t.Run("Chatflow_ConversationDel", func(t *testing.T) {
		err := client.ChatflowApp().ConversationDel("adc2ad24-fa4e-4dbb-8c16-ead1eaaa6c38", "")
		t.Logf("err=%v", err)
	})

	t.Run("Chatflow_ConversationRename", func(t *testing.T) {
		resp, err := client.ChatflowApp().ConversationRename(types.ConversationRenameReq{
			ConversationId: "f6da1bba-6341-42ed-9021-4a88b2f0dd0a",
			Name:           "修改后的新名称",
			AutoGenerate:   false,
			User:           "",
		})
		t.Logf("resp=%v err=%v", resp, err)
	})

	t.Run("Chatflow_AudioToText", func(t *testing.T) {
		client.ChatflowApp().AudioToText("", nil, "")
	})

	t.Run("Chatflow_TextToAudio", func(t *testing.T) {
		err := client.ChatflowApp().TextToAudio(types.Text2Audio{
			MessageId: "",
			Text:      "你是谁？今天是几号",
			User:      "",
		})
		t.Logf("resp=%v err=%v", "", err)
	})

	t.Run("Chatflow_AppMeta", func(t *testing.T) {
		resp, err := client.ChatflowApp().AppMeta()
		t.Logf("resp=%v err=%v", resp, err)
	})
}
