/**
 * Created by cdimonaco on 09/07/2017.
 */


const startQuizDialog = document.getElementById("startquiz");
const questionGroup = document.getElementById("questiongroup");
const resultForm = document.getElementById("result");
const responseArea = document.getElementById("backendresponse");
const scoreBoardArea = document.getElementById("scoreboard");
const apiUrl = "http://mockbin.org/bin/";

let playerUsername = "";
let quizQuestions = {};
let idQuiz = "";


function handleApiErrors(response){
    console.log(response);
    if(!response.ok){
        let error = new Error("Error");
        error.reason = response.json();
        throw error;
    }

    return response;
}


function startQuiz() {
    startQuizDialog.setAttribute("style","display:none;");
    playerUsername = document.getElementById("startusername").value;
    populateQuiz();
}


function callApi(url,body,method,success){
    console.log(url);
    let headers = new Headers();
    headers.append("Content-Type","application/json");
    //headers.append('Authorization', 'Basic ' + btoa(playerUsername + ":" + "x"));

    let requestOpt = {
      redirect:"follow",
        headers:headers,
        mode:"cors"
    };
    requestOpt["method"] = method;
    method === "POST" ? requestOpt["body"] = JSON.stringify(body) :false;
    let request = new Request(apiUrl + url,requestOpt);

    fetch(request)
        .then(handleApiErrors)
        .then(function (response) {
            console.log(response);
            return response.json();
        })
        .then(function (data) {
            console.log(data);
            return success(data);
        })
        .catch(function (error) {
            console.log(error);
        })
}


function populateQuiz() {
    callApi("623dc115-9d9c-4a5c-8c10-3a3ca4aa6da6",{},"GET", function (data) {
       console.log("Api call successful",data) ;
       quizQuestions = data.questions;
       idQuiz=data.id;
        questionGroup.setAttribute("style","display:inherit,margin-top:60px");
        console.log("POPULATE QUIZ",data,data.questions.length);
        for(i=0;i<data["questions"].length;i++){
            console.log("into the loop")
            let newQuestionDiv = document.createElement("div");
            newQuestionDiv.id = "question"+i;
            newQuestionDiv.classList.add("question")
            let content = `<div class="panel panel-default">
                    <div class="panel-heading">
                        Domanda ${i+1}
                    </div>
                    <div class="panel-body">
                        <div class="alert alert-info">
                            ${data.questions[i].question}
                        </div>
                        <form name="${"question"+i}">
                            <div class="form-group">`
            for(j=0;j<data.questions[i].answers.length;j++){
                content += `<div class="radio"><label><input name="answer" type="radio" value="1">${data.questions[i].answers[i]}</label></div>`
            }
            content += `<div class="radio"><label><input name="answer" type="radio" value="-1" checked="true">Salta</label></div>
                            </div>
                        </form>
                    </div>
                </div>` ;
            newQuestionDiv.innerHTML = content;
            questionGroup.appendChild(newQuestionDiv);
        }
        resultForm.setAttribute("style","display:inherit;")
    });
}


function buildAnswers(){
    let answer = [];
    //Loop through questions, and pick up the selected answer
    for(i=0;i<quizQuestions.length;i++){
        let currentAnswer = parseInt(document.forms["question"+i]["answer"].value);
        answer.push(currentAnswer);
    }
    return answer;
}

function sendResults() {
    let quizAnswers = buildAnswers();
    let backendRequest = {
        answers:quizAnswers,
    };
    console.log("Send results");
    console.log("Send to backend:",JSON.stringify(backendRequest));
    callApi("1efffd90-e903-41dd-b511-efd0bbe4ce34",backendRequest,"POST", function (data) {
        receiveResponse(data.results,data.score);
    });
}


function receiveResponse(data,score){
    questionGroup.setAttribute("style","display:none;");
    resultForm.setAttribute("style","display:none");
    console.log("Populating response");
     for(i=0;i<data.length-1;i++){
         console.log("Indice i",quizQuestions[i].answers);
            let newAnswerDiv = document.createElement("div");
            newAnswerDiv.id = "question"+i;
            newAnswerDiv.classList.add("question")
            newAnswerDiv.innerHTML = `<div class="panel panel-default">
                    <div class="panel-heading">
                        Domanda ${i+1}
                    </div>
                    <div class="panel-body">
                        <div class="alert alert-info">
                            ${quizQuestions[i].question}
                        </div>
                            <p><strong>Hai risposto: ${quizQuestions[i].answers[data[i].given]}</strong></p>
                            <p><strong>La risposta corretta è:  ${quizQuestions[i].answers[data[i].correct]}</strong></p>
                    </div>
                </div>`
            responseArea.appendChild(newAnswerDiv);
        }
    let scoreDiv = `<div class="alert alert-info"><p>Il tuo score è : ${score}</p></div>`
    responseArea.innerHTML += scoreDiv;
    responseArea.setAttribute("style","display:inherit");
}



function getScores() {
    callApi("7a4feb7a-ca2a-4a02-b5e3-d0042aab9abb",{},"GET", function (data) {
        let scoreTableBody = document.querySelector("#scoreboard > table");
        responseArea.setAttribute("style", "display:none");
        for (i = 0; i < data.scores.length; i++) {
            let newScoreRow = scoreTableBody.insertRow();
            let cellName = newScoreRow.insertCell(0);
            let cellScore = newScoreRow.insertCell(0);
            cellName.innerHTML = ` ${data.scores[i].score}`;
            cellScore.innerHTML = `${data.scores[i].user}`;
        }
        scoreBoardArea.setAttribute("style", "display:inherit");
        console.log("Score");
    });
}