// pages/chat/chat.js
const app = getApp()

Page({
  data: {
    messages: [], // 存储聊天消息
    inputValue: '', // 输入框内容
    toView: '' // 滚动到指定消息
  },

  onLoad() {
    // 页面加载时可以加载历史聊天记录
    this.loadChatHistory();
  },

  loadChatHistory() {
    // TODO: 从后端获取历史聊天记录
    // 示例：
    // wx.request({
    //   url: `${app.globalData.backendUrl}/api/v1/chat/history`,
    //   success: (res) => {
    //     this.setData({
    //       messages: res.data.history
    //     });
    //     this.scrollToBottom();
    //   }
    // });
  },

  onInput(e) {
    this.setData({
      inputValue: e.detail.value
    })
  },

  sendMessage() {
    const prompt = this.data.inputValue.trim()
    if (!prompt) {
      return
    }

    const newMessages = [...this.data.messages, { role: 'user', content: prompt }]
    this.setData({
      messages: newMessages,
      inputValue: '',
      toView: `msg-${newMessages.length - 1}` // 滚动到最新消息
    })

    wx.request({
      url: `${app.globalData.backendUrl}/api/v1/gemini/generate`, // 使用全局配置的后端API地址
      method: 'POST',
      timeout: 30000, // 增加超时时间到30秒 (30000毫秒)
      header: {
        'Content-Type': 'application/json'
      },
      data: {
        messages: newMessages // 发送完整的消息列表以支持多轮对话
      },
      success: (res) => {
        const aiResponse = res.data.content || '抱歉，AI未能生成内容。'
        const updatedMessages = [...this.data.messages, { role: 'ai', content: aiResponse }]
        this.setData({
          messages: updatedMessages,
          toView: `msg-${updatedMessages.length - 1}`
        })
      },
      fail: (err) => {
        console.error('请求失败', err)
        const errorMessage = '请求AI服务失败，请稍后再试。'
        const updatedMessages = [...this.data.messages, { role: 'ai', content: errorMessage }]
        this.setData({
          messages: updatedMessages,
          toView: `msg-${updatedMessages.length - 1}`
        })
      }
    })
  },

  scrollToBottom() {
    // 确保滚动到最新消息
    wx.createSelectorQuery().select('.message-list').boundingClientRect((rect) => {
      if (rect && rect.height) {
        this.setData({
          toView: `msg-${this.data.messages.length - 1}`
        });
      }
    }).exec();
  }
})