import type { Email } from 'types'
import { Dispatch, MouseEventHandler, ReactElement, SetStateAction } from 'react'
import { formatDistanceToNow, parseJSON } from 'date-fns'

type EmailsTableProps = {
  emails: Email[]
  setSelected: Dispatch<SetStateAction<Email | null>>
}

const EmailsTable = (props: EmailsTableProps): ReactElement => {
  const { emails, setSelected } = props;

  if (emails.length === 0) {
    return (
      <div className="prose">
        <p>No emails have been received yet.</p>
      </div>
    )
  }

  const rendered = emails.map((email: Email, idx: number) => 
    <EmailRow idx={idx} key={idx} email={email} onClick={() => setSelected(email)} />
  )

  return (
    <div className="overflow-auto border border-gray-200 sm:rounded-lg">
      <table className="min-w-full divide-y divide-gray-200 table-fixed">
        <thead className="bg-gray-50">
          <tr className="border-b">
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">From</th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Subject</th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Received</th>
          </tr>
        </thead>

        <tbody className="bg-white divide-y divide-gray-200">
          {rendered}
        </tbody>
      </table>
    </div>
  )
}

const EmailRow = ({ email, onClick, idx }: {email: Email, onClick: MouseEventHandler, idx: number }): ReactElement => {
  const receivedAt: Date = parseJSON(email.received_at);
  return (
    <tr data-test={`email-row-${idx}`} onClick={onClick}>
      <td className="px-6 py-4 whitespace-nowrap">{email.from}</td>
      <td data-test={`email-row-subject`} className="px-6 py-4 whitespace-nowrap overflow-ellipsis overflow-hidden">{email.subject}</td>
      <td className="px-6 py-4 whitespace-nowrap">{formatDistanceToNow(receivedAt)} ago</td>
    </tr>
  )
}

export default EmailsTable
