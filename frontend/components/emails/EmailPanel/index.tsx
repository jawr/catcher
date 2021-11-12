import { ReactElement } from "react";
import type { Email } from 'types'
import DOMPurify from 'dompurify'
import { format, parseJSON } from 'date-fns'

type EmailPanelProps = {
  email: Email
}

const EmailPanel = (props: EmailPanelProps): ReactElement => {
  const { email } = props;

  const receivedAt: Date = parseJSON(email.received_at);

  const text: string = Buffer.from(email.text, 'base64').toString()
  const html: string = Buffer.from(email.html, 'base64').toString()

  const message = (html.length > 0) ?
    <div dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(html) }} />
    :
    <pre className="whitespace-pre-line break-words">{text}</pre>

  return (
    <>
      <div>
        <h2 className="text-2xl">{email.subject}</h2>
        <h3 className="font-bold">{email.from}</h3>
        <p>{format(receivedAt, "yyyy-MM-dd hh:mm:ss")}</p>
      </div>

      <div className="mt-6">
        {message}
      </div>
      </>
  )
}

export default EmailPanel
