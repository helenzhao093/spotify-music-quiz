import React from 'react';

function TrackOptions({tracks, trackOptionIndexes, evaluateAnswer}) {
    return (
        <div>
        {
            trackOptionIndexes.map(index => {
                return <button id={tracks[index].ID} onClick={() => evaluateAnswer(index)}>{tracks[index].Name}</button>
            })
        }
        </div>
    )
}

export default TrackOptions;