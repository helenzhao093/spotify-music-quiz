import React from 'react'; // modules 

function AudioPlayer({track}) {
    //console.log('audio ' + track.Name)
    React.useEffect(() => {
        playVid(track);
    }, [track])
  
    function playVid(track) {
        //console.log('audio load new ' + track.Name);
        let vid = document.getElementById("videoPlayer");
        let source = document.getElementById("musicsource");
        source.setAttribute('src', track.PreviewURL);
        vid.load();
        vid.play();
    }
    return (
        <video id="videoPlayer" controls autoplay name="media">
            <source id="musicsource" src="" type="audio/mpeg"/>
        </video>
    )
} 

export default AudioPlayer;