import { useQuery } from '@tanstack/react-query'

const fetchAthletes = async () => {
  const res = await fetch('http://localhost:8080/api/athletes')
  if (!res.ok) throw new Error('Failed to fetch athletes')
  return res.json()
}

function AthleteList() {
  const { data: athletes, isLoading, isError, error } = useQuery({
    queryKey: ['athletes'],
    queryFn: fetchAthletes,
  })

  if (isLoading) {
    return <p className="text-center text-gray-600">Loading athletes...</p>
  }

  if (isError) {
    return <p className="text-center text-red-600">Error: {error.message}</p>
  }

  return (
    <div className="max-w-2xl mx-auto">
      <table className="w-full bg-white rounded-lg shadow-md overflow-hidden">
        <thead className="bg-blue-500 text-white">
          <tr>
            <th className="py-3 px-4 text-left">Name</th>
            <th className="py-3 px-4 text-left">Grade</th>
            <th className="py-3 px-4 text-left">Personal Record</th>
          </tr>
        </thead>
        <tbody>
          {athletes.map((athlete) => (
            <tr key={athlete.id} className={athlete.id % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>
              <td className="py-3 px-4">{athlete.name}</td>
              <td className="py-3 px-4">{athlete.grade}</td>
              <td className="py-3 px-4">{athlete.personalRecord}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default AthleteList
