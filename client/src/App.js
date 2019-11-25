import React from 'react'; // modules 
import logo from './logo.svg';
import './App.css'; 
import './PlayListOptions.js'
import './TrackOptions.js'

function App() {
  const [playlists, setPlaylists] = React.useState([]);
  const [tracks, setTracks] = React.useState([]);
  const [showTracks, setShowTracks] = React.useState(false);
  const [currentTrackIndex, setCurrentTrackIndex] = React.useState({});
  const [trackOptionIndexes, setTrackOptionIndexes] = React.useState([]);
  const numberOptions = 5;
  const [score, setScore] = React.useState(0);
  const [timeoutFn, setTimeoutFn] = React.useState(() => {});
  const [songTimeoutFn, setSongTimeoutFn] = React.useState(() => {});

  React.useEffect(() => {
      getPlaylists();
  }, [])

  const getPlaylists = async () => {
      console.log('getting playlists');
      /*fetch( "http://localhost:8080/getPlaylists")
          .then(data => {
              return data.clone().json() 
          }).then(response => {
              console.log(response);
              setPlaylists(response);
          }) */
          let response = await fetch( "http://localhost:8080/getPlaylists");
          response = await response.json();
          console.log(response);
          setPlaylists(response);
      }

  
  function getTracks(playlistId) {
      console.log('getting tracks ' + playlistId)
      
      fetch( "http://localhost:8080/getTracksFromPlaylist/" + playlistId)
          .then( data => {
              return data.clone().json();
          }).then(response => {
              console.log(response);
              setTracks(response);
              setRound(response);
              setShowTracks(true);
          })
  }

  function evaluateAnswer(selectedIndex) {
      if (selectedIndex === currentTrackIndex) {
          setScore(score + 10)
      }
      setRound(tracks)
  }

  function setRound(tracks) {
      console.log(timeoutFn);
      window.clearTimeout(timeoutFn); // clear all the timeouts
      window.clearTimeout(songTimeoutFn);
      setTimeoutFn(() => {});
      setSongTimeoutFn(() => {});
      let trackIndex = Math.floor(Math.random() * (tracks.length - 1));
      while (currentTrackIndex && trackIndex == currentTrackIndex) {
          console.log('finding track to play')
          trackIndex = Math.floor(Math.random() * (tracks.length - 1));
      }
      let options = new Set([trackIndex]);
      while (options.size < numberOptions) {
          options.add(Math.floor(Math.random() * (tracks.length - 1)));
      }
      let optionsArr = [...options];
      shuffleArray(optionsArr);
      setCurrentTrackIndex(trackIndex);
      setTrackOptionIndexes(optionsArr);
      setTimeoutForSong(tracks);
      removeOptionTimer(optionsArr, trackIndex);
  }

  function setTimeoutForSong(tracks) {
      setSongTimeoutFn(setTimeout(function() {
          setRound(tracks);
      }, 30000));
  }

  function removeOptionTimer(optionsArr, trackIndex) {
      if (optionsArr.length > 2) {
          setTimeoutFn(setTimeout(function () {
              let a = optionsArr.splice(0)
              let index = Math.floor(Math.random()*a.length);
              while (a[index] == trackIndex) {
                  console.log(a.length);
                  console.log(index);
                  console.log('finding index to remove');
                  index = Math.floor(Math.random()*a.length);
              }
              console.log(index);
              a.splice(index, 1)
              console.log(a);
              setTrackOptionIndexes(a);
              removeOptionTimer(a, trackIndex);
          }, 5000));
          
      }
  }

  function shuffleArray(array) {
      for (var i = array.length - 1; i > 0; i--) {
          var j = Math.floor(Math.random() * (i + 1));
          var temp = array[i];
          array[i] = array[j];
          array[j] = temp;
      }
  }

  if (showTracks) {
      return (
          <div>
              <p>{score}</p>
              <TrackOptions tracks={tracks} trackOptionIndexes={trackOptionIndexes} evaluateAnswer={evaluateAnswer}/>
              <AudioPlayer track={tracks[currentTrackIndex]} />
          </div>
      )
  } else {
      return (
          <PlayListOptions playlists={playlists} getTracks={getTracks}/>
      )
  } 
}

/*function PlayListOptions({playlists, getTracks}) {
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
} */

export default App;
