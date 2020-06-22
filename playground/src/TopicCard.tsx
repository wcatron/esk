import { Client, Message, MessageCommand } from "esk-client-typescript";
import React, { useState, useEffect, useMemo } from "react";

export const TopicCard:React.FC<{
    topic: string,
    onPublish: (payload: string) => void,
    client: Client
}> = ({
    topic,
    onPublish,
    client
}) => {
    const [payload, setPayload] = useState('')
    const [messageCount, setMessageCount] = useState(0)
    const [messages, setMessages] = useState<Message[]>([])
    useMemo(() => {
        const existingState = localStorage.getItem(topic) || "[]"
        try {
            const state = JSON.parse(existingState) as {
                command: MessageCommand,
                data: Record<number, number>
            }[]
            const messages = state.map(_element => {
                const values = Object.values(_element.data)
                const message = new Message({
                    command: _element.command,
                    topic,
                    data: new Uint8Array(values)
                })
                return message
            })
            setMessages(messages)
        } catch (e) {
            console.error(e)
        } finally {
            return []
        }
    }, [topic, messageCount, setMessages])

    useEffect(() => {
        const callback = (message: Message) => {
            if (message.topic && 
                (message.command === MessageCommand.PUBLISH || message.command === MessageCommand.INFORM)) {
                const existingState = localStorage.getItem(topic) || "[]"
                try {
                    const state = JSON.parse(existingState) as any[]
                    state.push(message)
                    localStorage.setItem(topic, JSON.stringify(state))
                    setMessageCount(state.length)
                } catch (e) {
                    console.error(e)
                }
            }
        }
        localStorage.removeItem(topic)
        client.subscribe(topic, callback)
        return () => {
            client.unsubscribe(topic, callback)
        }
    }, [topic, client, setMessageCount])

    return <div className='TopicCard'>
        <h3>Topic: {topic}</h3>
        <ul>
            {messages.map((message, index) => <li key={index}>{JSON.stringify(message.payload)}</li>)}
        </ul>
        <input value={payload} onChange={e => setPayload(e.currentTarget.value)} /><button onClick={() => {
            onPublish(payload)
        }}>Publish</button>
    </div>
}