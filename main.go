package main

import (
    "github.com/lxn/walk"
    . "github.com/lxn/walk/declarative"
    "github.com/lxn/win"
    "simulate-client/client"
)

var (
    mainWidth  int32 = 800
    mainHeight int32 = 500
)

func main() {
    // 输入 输出框
    var inTE, outTE *walk.TextEdit

    // 连接按钮
    var connectButton *walk.PushButton
    // 连接地址输入框
    var wssUrl *walk.LineEdit
    // 主窗口
    var mainWindow *walk.MainWindow

    var channel *client.Channel
    MainWindow{
        AssignTo: &mainWindow,
        Title:    "消息测试",
        MinSize:  Size{300, 300},
        Size:     Size{Width: int(mainWidth), Height: int(mainHeight)},
        Layout:   VBox{MarginsZero: true, Spacing: 10},
        Children: []Widget{
            VSplitter{
                Children: []Widget{
                    Composite{
                        Layout: HBox{
                            Spacing: 5,
                            Margins: Margins{5, 5, 0, 5},
                        },
                        Children: []Widget{
                            LinkLabel{
                                Font: Font{
                                    Bold: true,
                                },
                                Text: "服务器url:",
                            },
                            LineEdit{
                                MinSize:     Size{30, 20},
                                MaxSize:     Size{500, 20},
                                Row:         200,
                                ToolTipText: "websocket连接地址",
                                AssignTo:    &wssUrl,
                            },
                            PushButton{
                                MaxSize:  Size{80, 20},
                                Text:     "连接",
                                AssignTo: &connectButton,
                                OnClicked: func() {
                                    if connectButton.Text() == "连接" {
                                        connectButton.SetText("断开")
                                        println(wssUrl.Text())
                                        channel = client.GetClientChannel(wssUrl.Text())
                                        outTE.AppendText("连接服务成功....\r\n")
                                        channel.Read()
                                        // 注册接收消息处理函数
                                        registerReceive(channel, outTE)

                                        // 当连接断开时 重置按钮文本
                                        channel.RegisterClose(func() {
                                            connectButton.SetText("连接")
                                        })

                                    } else {
                                        connectButton.SetText("连接")
                                        outTE.AppendText("断开服务成功....\r\n")
                                        if channel != nil {
                                            channel.Close()
                                        }
                                    }
                                },
                            },
                        },
                    },
                    HSplitter{
                        Children: []Widget{
                            TextEdit{AssignTo: &inTE, VScroll: true, Font: Font{PointSize: 12}, RowSpan: 20, Row: 20},
                            TextEdit{AssignTo: &outTE, ReadOnly: true, VScroll: true, RowSpan: 50},
                        },
                    },
                },
            },
            Composite{
                Layout: HBox{},
                Children: []Widget{
                    PushButton{
                        MaxSize: Size{300, 60},
                        Text:    "清空",
                        Font: Font{
                            PointSize: 11,
                        },
                        OnClicked: func() {
                            outTE.SetText("")
                        },
                    },
                    PushButton{
                        MaxSize: Size{300, 60},
                        Text:    "发送",
                        Font: Font{
                            PointSize: 11,
                        },
                        OnClicked: func() {
                            if channel == nil {
                                outTE.AppendText("还未连接到服务器，请先连接...\r\n")
                                return
                            }
                            if inTE.Text() == "" {
                                return
                            }
                            outTE.AppendText(" >>>>>>> : " + inTE.Text() + "\r\n")
                            channel.SendMsg(inTE.Text())
                            //inTE.SetText("")
                        },
                    },
                },
            },
        },
    }.Create()

    xScreen := win.GetSystemMetrics(win.SM_CXSCREEN)
    yScreen := win.GetSystemMetrics(win.SM_CYSCREEN)
    win.SetWindowPos(mainWindow.Handle(), 0, (xScreen-mainWidth)/2, (yScreen-mainHeight)/2, mainWidth, mainHeight, win.SWP_FRAMECHANGED)
    win.ShowWindow(mainWindow.Handle(), win.SW_SHOW)
    mainWindow.Run()
}

func registerReceive(channel *client.Channel, edit *walk.TextEdit) {
    go func() {
        for {
            select {
            case <-channel.GetCloseChannel():
                return
            default:
                msg := channel.GetMsg()
                edit.AppendText(" <<<<<<<< : " + msg + "\r\n")
            }
        }
    }()
}
