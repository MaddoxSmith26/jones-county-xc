import { useState } from 'react'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col items-center justify-center">
      <h1 className="text-4xl font-bold text-gray-800 mb-8">
        Hello World
      </h1>
      <div className="bg-white p-8 rounded-lg shadow-md">
        <p className="text-gray-600 mb-4">Welcome to the application</p>
        <button
          onClick={() => setCount((count) => count + 1)}
          className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded"
        >
          Count: {count}
        </button>
      </div>
    </div>
  )
}

export default App
