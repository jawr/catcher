import type { NextPage } from 'next' 
import type { Email } from 'types'
import { useRouter } from 'next/router'
import { useEffect, useState } from 'react'
import Nav from 'components/Nav'
import EmailsTable from 'components/emails/Table'
import EmailPanel from 'components/emails/EmailPanel'


const Emails: NextPage = () => {
  const router = useRouter()
  const key: string = router.query["key"]?.toString().toLowerCase() ?? ''
  const [selected, setSelected] = useState<Email | null>(null)
  const [emails, setEmails] = useState<Email[]>([])
  const [error, setError] = useState<Error | null>(null)


  useEffect(() => {
    if (!key) {
      return
    }

    const ws = new WebSocket(`${process.env.NEXT_PUBLIC_WS_URI}/api/v1/subscribe`, 'binary')

    ws.onopen = () => {
      ws.send(JSON.stringify({key}))
    }

    ws.onerror = (event: Event) => {
      setError(new Error("error reading data from websocket: " + event))
    }

    ws.onmessage = (event: MessageEvent) => {
      try {
          const emails: Email[] = JSON.parse(event.data);
          setEmails(emails);
      } catch (error) {
        setError(error as Error)
      }
    }
  },[key])

  return (
    <>
      <Nav />
      <div className="xl:px-24 container mx-auto mt-24 md:mx-auto px-5 md:px-0">
        <div className="">
          <h1 className="text-4xl font-thin">Inbox</h1>
          <p className="tracking-wide font-light">{key}@catcher.mx.ax</p>
        </div>

        <div className="mt-6">
          {error && <p>An unexpected error occured: {error.message}</p>}
          {emails && <EmailsTable emails={emails} setSelected={setSelected} />}
        </div>

        {selected && <div className="mt-6 border p-6 border-b border-gray-200 sm:rounded-lg">
          <EmailPanel email={selected} />
        </div>}
      </div>
      </>
  )
}

export default Emails
