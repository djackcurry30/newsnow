// This file is used to trigger opening the independent window
// It immediately closes the popup and opens the window

import { useEffect } from "react"

function PopupTrigger() {
  useEffect(() => {
    // Open independent window
    chrome.windows.create({
      url: chrome.runtime.getURL("tabs/window.html"),
      type: "popup",
      width: 620,
      height: 540,
      left: Math.round((window.screen.width - 620) / 2),
      top: Math.round((window.screen.height - 540) / 2),
      focused: true
    })

    // Close this popup
    window.close()
  }, [])

  return (
    <div style={{
      width: "200px",
      height: "100px",
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
      fontSize: "14px",
      color: "#666"
    }}>
      正在打开窗口...
    </div>
  )
}

export default PopupTrigger
