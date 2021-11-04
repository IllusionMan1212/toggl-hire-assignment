import { useState, useEffect } from "react";

export default function TextQuestion({setAnswers, ...props}) {
  const [answer, setAnswer] = useState({
    questionId: props.id,
    enteredText: "",
  });

  const handleChange = (e) => {
    const text = e.target.value;

    setAnswer({
      ...answer,
      enteredText: text
    });
  };

  useEffect(() => {
    setAnswers(answers => {
      return answers.map(_answer => {
        if (_answer.questionId === answer.questionId) {
          _answer = {
            ..._answer,
            enteredText: answer.enteredText
          };
        }
        return _answer;
      });
    });
  }, [answer, setAnswers]);

  return (
    <div className="question">
      <h5>Question #{props.number} - Free Text</h5>
      <p>{props.body}</p>
      <textarea onChange={handleChange}/>
    </div>
  )
}
