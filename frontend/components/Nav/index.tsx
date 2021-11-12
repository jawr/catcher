import React, { ReactElement } from "react"
import Head from 'next/head'
import Link from 'next/link'

const Nav: React.FC = (): ReactElement => {
  return (
    <>
      <Head>
        <title>Catcher</title>
        <meta name="description" content="Temporary email service to enable QA and automated testing in software development." />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <header className="text-gray-600 body-font mb-6 absolute left-0 right-0 top-0">
        <div className="container mx-auto flex flex-wrap py-5 flex-col md:flex-row items-center">
          <Link href="/"><span className="text-xl cursor-pointer">Catcher</span></Link>
          <nav className="md:ml-auto flex flex-wrap items-center text-base flex-end">
            <a href="https://github.com/jawr/catcher" className="ml-5 hover:text-gray-900">GitHub</a>
          </nav>
        </div>
      </header>
    </>
  )
}

export default Nav
