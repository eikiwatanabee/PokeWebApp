interface Props {
  status: string
}

const statusConfig: Record<string, { label: string; className: string }> = {
  unread: { label: '未読', className: 'bg-gray-100 text-gray-700' },
  reading: { label: '読書中', className: 'bg-blue-100 text-blue-700' },
  finished: { label: '読了', className: 'bg-green-100 text-green-700' },
}

export function StatusBadge({ status }: Props) {
  const config = statusConfig[status] || { label: status, className: 'bg-gray-100 text-gray-700' }
  return (
    <span className={`inline-block px-2 py-1 rounded-full text-xs font-medium ${config.className}`}>
      {config.label}
    </span>
  )
}
