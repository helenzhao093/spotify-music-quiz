import React from 'react';

function PlayListOptions({playlists, getTracks}) {
    return (
        <div>
        {
            playlists.map(playlist => {
            return <button onClick={() => getTracks(playlist.ID)} id={playlist.ID}>{playlist.Name}</button>
            })
        }
        </div>
    )
}

export default PlayListOptions;