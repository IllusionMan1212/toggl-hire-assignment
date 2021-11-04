import Option from "./Option";
import { useState, useEffect } from "react";
import { sortByWeight } from "../utils";

export default function ChoiceQuestion({setAnswers, ...props}) {
  const [options, setOptions] = useState([]);
  const [answer, setAnswer] = useState({
    questionId: props.id,
    selectedOptionId: -1,
    selectedOptionWeight: -1
  });

  useEffect(() => {
    const sortedOptions = props.options.sort(sortByWeight);
    setOptions(sortedOptions);
  }, [props.options]);

  useEffect(() => {
    setAnswers(answers => {
      return answers.map(_answer => {
        if (_answer.questionId === answer.questionId) {
          _answer = {
            ..._answer,
            selectedOptionId: answer.selectedOptionId,
            selectedOptionWeight: answer.selectedOptionWeight
          }
        }
        return _answer;
      });
    });
  }, [answer, setAnswers]);

  return (
    <div className="question">
      <h5>Question #{props.number} - Multiple Choice</h5>
      <p>{props.body}</p>
      {options.map((option) => {
        return (
          <Option
            key={option.id}
            type={option.__typename}
            body={option.body}
            id={option.id}
            weight={option.weight}
            onClick={() => setAnswer({...answer, selectedOptionId: option.id, selectedOptionWeight: option.weight})}
          />
        );
      })}
    </div>
  )
}
