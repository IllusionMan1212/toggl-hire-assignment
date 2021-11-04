export default function Option(props) {
  return (
    <div className="option" onClick={props.onClick}>
      <label>
        <input type="radio" name="answer"/>
        {props.body}
      </label>
    </div>
  )
}
