// Background service worker - handles extension icon click to open independent window

chrome.action.onClicked.addListener(async () => {
  // Check if window already exists
  const windows = await chrome.windows.getAll({ populate: true })
  const existingWindow = windows.find(w =>
    w.tabs?.some(tab => tab.url?.includes("tabs/window.html"))
  )

  if (existingWindow) {
    // Focus existing window
    await chrome.windows.update(existingWindow.id!, { focused: true })
    return
  }

  // Get screen dimensions to center the window
  const screenWidth = window.screen.width
  const screenHeight = window.screen.height
  const width = 620
  const height = 540
  const left = Math.round((screenWidth - width) / 2)
  const top = Math.round((screenHeight - height) / 2)

  // Create new independent window
  await chrome.windows.create({
    url: chrome.runtime.getURL("tabs/window.html"),
    type: "popup",
    width,
    height,
    left,
    top,
    focused: true
  })
})

// Keep service worker alive
chrome.runtime.onStartup.addListener(() => {
  console.log("NewsNow extension started")
})

chrome.runtime.onInstalled.addListener(() => {
  console.log("NewsNow extension installed")
})
