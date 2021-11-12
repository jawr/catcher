export type Email = {
  from: string
  to: string
  subject: string
  headers: Record<string, string>
  html: string
  text: string
  received_at: string
}
