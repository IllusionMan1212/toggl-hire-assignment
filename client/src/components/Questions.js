import { gql, useQuery, useMutation } from "urql";
import ChoiceQuestion from "./ChoiceQuestion";
import TextQuestion from "./TextQuestion";
import { useState, useEffect } from "react";
import { sortByWeight } from "../utils";
import { useToastContext } from "../contexts/toastContext";
import { useNavigate } from "react-router-dom";

export function Questions() {
  const [{ data, fetching, error }] = useQuery({ query });
  const [questions, setQuestions] = useState([]);
  const [answers, setAnswers] = useState([]);
  const [, submitAnswers] = useMutation(sendAnswers);
  const toast = useToastContext();
  const navigate = useNavigate();

  useEffect(() => {
    if (data) {
      data.questions.sort(sortByWeight);
      setQuestions(data.questions);

      let initialAnswers = [];

      data.questions.forEach((question) => {
        initialAnswers = initialAnswers.concat({
            questionId: question.id,
            questionType: question.__typename,
            questionWeight: question.weight,
            enteredText: null,
            selectedOptionId: null,
            selectedOptionWeight: null
        });
      });
      setAnswers(initialAnswers);
    }
  }, [data]);

  const handleSubmit = () => {
    submitAnswers({ "answers": answers }).then(result => {
      if (result.error) {
        toast(`ERROR: ${result.error.message}`, 4000);
      } else {
        toast(result.data.submitAnswers.content, 3500);
        navigate("/results");
      }
    });
  };

  if (fetching) return "Loading...";
  if (error) return `Error: ${error}`;

  return (
    <div className="container">
      {questions.map((question, index) => {
        switch (question.__typename) {
        case "ChoiceQuestion":
            return (
              <ChoiceQuestion
                key={question.id}
                body={question.body}
                id={question.id}
                type={question.__typename}
                weight={question.weight}
                options={question.options}
                number={index + 1}
                setAnswers={setAnswers}
              />
            );
        case "TextQuestion":
            return (
              <TextQuestion
                key={question.id}
                body={question.body}
                id={question.id}
                weight={question.weight}
                options={question.options}
                number={index + 1}
                setAnswers={setAnswers}
              />
            );
        default:
            return (
              <div>Invalid question type</div>
            );
        }
      })}

      <div>
        <button onClick={handleSubmit}>Submit your answers</button>
      </div>
    </div>
  );
}

const sendAnswers = gql`
  mutation($answers: [AnswerInput!]!) {
    submitAnswers(answers: $answers) {
      content
    }
  }
`;

const query = gql`
  query {
    questions {
      __typename
      ... on ChoiceQuestion {
        id
        body
        weight
        options {
          id
          body
          weight
        }
      }
      ... on TextQuestion {
        id
        body
        weight
      }
    }
  }
`;
