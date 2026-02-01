import { useCallback, useEffect, useState } from "react"

import { api } from "~/lib/api"

import { Header } from "~/components/layout/Header"
import { Footer } from "~/components/layout/Footer"
import { ErrorBanner } from "~/components/layout/ErrorBanner"
import { ContentArea } from "~/components/layout/ContentArea"

import { COLUMNS } from "~/components/constants"
import type { SourceResponse } from "~/components/types"
import { getSourceName } from "~/components/utils"

import "~/globals.css"

function IndexPopup() {
  const [activeColumn, setActiveColumn] = useState("focus")
  const [sourceData, setSourceData] = useState<Record<string, SourceResponse>>({})
  const [loading, setLoading] = useState<Record<string, boolean>>({})
  const [error, setError] = useState("")

  // Initialize API client
  useEffect(() => {
    api.init()
  }, [])

  // Load news for a specific source
  const loadSource = useCallback(async (sourceId: string) => {
    setLoading((prev) => ({ ...prev, [sourceId]: true }))
    setError("")

    try {
      const data = await api.request<SourceResponse>(`/api/s/${sourceId}`)
      setSourceData((prev) => ({ ...prev, [sourceId]: data }))
    } catch (err: any) {
      console.error(`Failed to load ${sourceId}:`, err)
      setError(`加载 ${getSourceName(sourceId)} 失败`)
    } finally {
      setLoading((prev) => ({ ...prev, [sourceId]: false }))
    }
  }, [])

  // Load all sources in active column
  const loadColumnNews = useCallback(
    async (columnId: string) => {
      const column = COLUMNS.find((c) => c.id === columnId)
      if (!column) return

      // Load each source in parallel
      await Promise.all(column.sources.map((sourceId) => loadSource(sourceId)))
    },
    [loadSource]
  )

  // Load initial data
  useEffect(() => {
    loadColumnNews("focus")
  }, [loadColumnNews])

  // Refresh all sources in active column
  const refreshAll = useCallback(async () => {
    setError("")
    await loadColumnNews(activeColumn)
  }, [activeColumn, loadColumnNews])

  // Handle column change
  const handleColumnChange = (columnId: string) => {
    setActiveColumn(columnId)
    loadColumnNews(columnId)
  }

  const handleNewsItemClick = (url: string) => {
    if (url && url !== "#") {
      window.open(url, "_blank")
    }
  }

  const isAnyLoading = Object.values(loading).some(Boolean)

  return (
    <div className="w-full h-full min-w-[600px] min-h-[500px] bg-gray-100 flex flex-col overflow-hidden">
      <Header
        activeColumn={activeColumn}
        isLoading={isAnyLoading}
        onColumnChange={handleColumnChange}
        onRefresh={refreshAll}
      />

      {error && <ErrorBanner message={error} />}

      <ContentArea
        activeColumn={activeColumn}
        sourceData={sourceData}
        loading={loading}
        onRefresh={loadSource}
        onItemClick={handleNewsItemClick}
      />

      <Footer />
    </div>
  )
}

export default IndexPopup
