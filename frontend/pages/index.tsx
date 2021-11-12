import type { NextPage } from 'next'
import Router from 'next/router'
import Nav from 'components/Nav'
import Search from 'components/Search'

type RandomKeyResponse = {
  key: string,
}

const Home: NextPage = () => {
  const handleCreate = async () => {
    const response: Response = await fetch(`/api/v1/random`)
    const random: RandomKeyResponse = await response.json()
    Router.push(`/email/${random.key}`)
  }

  return (
    <>
      <Nav />

      <div className="grid w-full min-h-screen place-content-center">
        <div className="mb-8">
          <h1 className="text-6xl font-thin">Check Inbox</h1>
        </div>

        <Search />

        <div className="mt-8 text-center tracking-wide font-light">
          <p>enter your email key or <span className="cursor-pointer hover:underline text-blue-400" onClick={handleCreate}>create</span> a new one</p>
        </div>

      </div>
    </>
  )
}

export default Home
