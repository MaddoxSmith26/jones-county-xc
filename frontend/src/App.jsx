import AthleteList from './components/AthleteList'

function App() {
  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <h1 className="text-4xl font-bold text-gray-800 mb-8 text-center">
        Jones County XC
      </h1>
      <h2 className="text-2xl font-semibold text-gray-700 mb-4 text-center">
        Athletes
      </h2>
      <AthleteList />
    </div>
  )
}

export default App
