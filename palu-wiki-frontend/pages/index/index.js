// index.js
// 获取应用实例
const app = getApp()

Page({
  data: {
    motto: 'Hello Palu Wiki',
    userInfo: {},
    hasUserInfo: false,
    canIUse: wx.canIUse('button.open-type.getUserInfo'),
    canIUseGetUserProfile: false,
    canIUseOpenData: wx.canIUse('open-data'), // 如需尝试获取用户信息可改为false
    messages: [], // 存储聊天消息
    inputValue: '', // 输入框内容
    toView: '' // 滚动到指定消息
  },
  // 事件处理函数
  bindViewTap() {
    wx.navigateTo({
      url: '../logs/logs'
    })
  },
  onLoad() {
    if (wx.getUserProfile) {
      this.setData({
        canIUseGetUserProfile: true
      })
    }
  },
  getUserProfile(e) {
    // 推荐使用wx.getUserProfile获取用户信息，开发者每次通过该接口获取用户个人信息均需用户确认
    // 开发者妥善保管用户快速登录后的 openid，以便后续无需弹窗二次确认
    wx.getUserProfile({
      desc: '用于完善会员资料', // 声明获取用户个人信息后的用途，后续会展示在弹窗中，请谨慎填写
      success: (res) => {
        this.setData({
          userInfo: res.userInfo,
          hasUserInfo: true
        })
      }
    })
  },
  getUserInfo(e) {
    // 不推荐使用getUserInfo获取用户信息，直接使用button open-type="getUserInfo"
    console.log(e)
    this.setData({
      userInfo: e.detail.userInfo,
      hasUserInfo: true
    })
  },

  goToGuideList() {
    wx.navigateTo({
      url: '/pages/guide-list/guide-list'
    })
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
        prompt: prompt
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
  }
})