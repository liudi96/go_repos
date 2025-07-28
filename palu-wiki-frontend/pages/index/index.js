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
    messages: [], // 存储聊天消息 (已废弃，AI聊天功能将迁移到独立页面)
    inputValue: '', // 输入框内容 (已废弃)
    toView: '' // 滚动到指定消息 (已废弃)
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

  goToChat() {
    wx.navigateTo({
      url: '/pages/chat/chat'
    })
  }
})