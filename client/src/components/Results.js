import { gql, useQuery } from "urql";
import { useState, useEffect } from "react";

export default function Results() {
  const [{ data, fetching, error }] = useQuery({ query });
  const [results, setResults] = useState({
    totalQuestions: 0,
    correctAnswers: 0,
    score: 0
  });

  useEffect(() => {
    if (data) {
      setResults({
        totalQuestions: data.results.totalQuestions,
        correctAnswers: data.results.correctAnswers,
        score: data.results.score
      });
    }
  }, [data]);

  if (fetching) return <div>Loading...</div>;
  if (error) return <div>Error: {JSON.stringify(error, null, 2)}</div>;

  return (
    <div className="container">
      <h2>Here's your results!</h2>
      <h4>You answered <b>{results.correctAnswers} / {results.totalQuestions}</b> questions correctly.</h4>
      <h3>With a score of {results.score.toFixed(2)}%</h3>
    </div>
  )
}

const query = gql`
  query {
    results {
      totalQuestions
      correctAnswers
      score
    }
  }
`
;
