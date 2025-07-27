// pages/guide-list/guide-list.js
const app = getApp()

Page({
  data: {
    guides: [],
    loading: false,
    error: ''
  },

  onLoad() {
    this.fetchGuides()
  },

  fetchGuides() {
    this.setData({ loading: true, error: '' })
    wx.request({
      url: `${app.globalData.backendUrl}/api/v1/guides`,
      method: 'GET',
      success: (res) => {
        if (res.statusCode === 200) {
          this.setData({
            guides: res.data.data,
            loading: false
          })
        } else {
          this.setData({
            error: '获取攻略失败',
            loading: false
          })
        }
      },
      fail: (err) => {
        console.error('请求失败', err)
        this.setData({
          error: '网络请求失败，请稍后再试。',
          loading: false
        })
      }
    })
  },

  goToGuideDetail(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/guide-detail/guide-detail?id=${id}`
    })
  }
})