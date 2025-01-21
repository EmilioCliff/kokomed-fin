import ReactDOM from 'react-dom';

function Spinner() {
  return ReactDOM.createPortal(
    <div className="loadingSpinnerContainer">
      <div className="loadingSpinner"></div>
    </div>,
    document.body
  );
}

export default Spinner;
