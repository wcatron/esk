import { Client, Message, MessageCommand } from "esk-client-typescript";
import React, { useState, useEffect, useMemo, useRef } from "react";

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
    const ref = useRef<HTMLOListElement>(null)
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
            if (messageCount !== messages.length) {
                console.error('Message count did not match messages array.')
            }
            setMessages(messages)
        } catch (e) {
            console.error(e)
        } finally {
            return []
        }
    }, [topic, messageCount, setMessages])

    useEffect(() => {
        if (ref.current) {
            ref.current?.scrollTo({
                top: 0
            })
        }
    }, [ref, messages])

    useEffect(() => {
        const callback = (message: Message) => {
            if (message.topic && 
                (message.command === MessageCommand.INFORM)) {
                const existingState = localStorage.getItem(topic) || "[]"
                try {
                    const state = JSON.parse(existingState) as any[]
                    state.push(message)
                    localStorage.setItem(topic, JSON.stringify(state))
                    const nextCursor = message.cursor! + message.payload.length;
                    localStorage.setItem(topic + ":cursor", String(nextCursor))
                    setMessageCount(state.length)
                } catch (e) {
                    console.error(e)
                }
            }
        }
        const cursor = parseInt(localStorage.getItem(`${topic}:cursor`) || '0')
        client.subscribe(topic, cursor, callback)
        return () => {
            client.unsubscribe(topic, callback)
        }
    }, [topic, client, setMessageCount])

    return <div className='TopicCard'>
        <h3>Topic: {topic}</h3>
        <ol ref={ref}>
            {messages.map((message, index) => <li key={index}>{JSON.stringify(message.payload)}</li>)}
        </ol>
        <input value={payload} onChange={e => setPayload(e.currentTarget.value)} /><button onClick={() => {
            onPublish(payload)
            setPayload('')
        }}>Publish</button>
    </div>
}