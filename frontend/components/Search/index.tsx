import React, { ReactElement } from "react"

const Search: React.FC = (): ReactElement => {
  return (
    <div className="flex items-center max-w-md mx-auto bg-white rounded-full border">
      <div className="w-full">
        <input type="search" className="w-full px-4 py-1 text-gray-900 rounded-full focus:outline-none"
          placeholder="" />
      </div>
      <div>
        <button type="submit" className="flex items-center justify-center w-12 h-12 text-gray-100 rounded-full bg-gray-500">
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
          </svg>
        </button>
      </div>
    </div>
  )
}

export default Search
