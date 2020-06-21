import React, { useState } from "react"

export const HostSelector: React.FC<{
    onSelect:(host: string) => void,
    defaultValue: string,
    connected: boolean
}> = ({ onSelect, defaultValue, connected }) => {
    const [ visible, setVisible ] = useState(false)
    const [ value, setValue] = useState(defaultValue)
    if (!visible && connected) {
        return <><button onClick={() => {
            setVisible(true)
        }}>{connected ? 'Connected' : 'Edit Host'}</button></>
    }
    return <>
        <button onClick={() => {
            onSelect(value)
            setVisible(false)
        }}>Change</button>
        <input value={value} onChange={(e) => {
            setValue(e.currentTarget.value)
        }} />
        <button onClick={() => {
            setVisible(false)
        }}>Cancel</button>
    </>
}