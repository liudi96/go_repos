// pages/guide-detail/guide-detail.js
const app = getApp()

Page({
  data: {
    guide: null,
    loading: false,
    error: ''
  },

  onLoad(options) {
    const id = options.id
    if (id) {
      this.fetchGuideDetail(id)
    } else {
      this.setData({ error: '缺少攻略ID' })
    }
  },

  fetchGuideDetail(id) {
    this.setData({ loading: true, error: '' })
    wx.request({
      url: `${app.globalData.backendUrl}/api/v1/guides/${id}`,
      method: 'GET',
      success: (res) => {
        if (res.statusCode === 200) {
          // Convert plain text content to rich-text nodes
          const formattedContent = this.formatContentForRichText(res.data.data.content);
          this.setData({
            guide: {
              ...res.data.data,
              content: formattedContent
            },
            loading: false
          })
        } else {
          this.setData({
            error: '获取攻略详情失败',
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
      },
    
      })
    },

    // Helper function to format content for rich-text component
    formatContentForRichText(content) {
      // Simple conversion: replace newlines with <br/> for basic paragraph breaks
      // For more complex formatting (e.g., Markdown to HTML), you'd need a more robust parser
      const htmlContent = content.replace(/\n/g, '<br/>');
      return htmlContent;
    }
})