const APIButton = (p: {
  url: string
  endpoint: string
  disabled: boolean
  setValue: (value: React.SetStateAction<boolean>) => void
  setError: (value: React.SetStateAction<string | undefined>) => void
}) => {
  const onClick = async () => {
    const resp = await fetch(`${p.url}/${p.endpoint}`)
    if (resp.ok) {
      p.setValue(true)
    } else {
      p.setError("Failed to start")
    }
  }
  return (
    <div>
      <button type="button" onClick={onClick} disabled={p.disabled}>
        {" "}
        {p.endpoint}
      </button>
    </div>
  )
}

export default APIButton
