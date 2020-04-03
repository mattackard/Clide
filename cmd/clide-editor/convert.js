const fileInput = document.getElementById("fileInput");
const jsonOutput = document.getElementById("jsonOutput");
const scriptOutput = document.getElementById("scriptOutput");
const saveFile = document.getElementById("saveFile");
const refresh = document.getElementById("refresh");
const run = document.getElementById("run");

fileInput.addEventListener("change", handleFiles, false);
saveFile.addEventListener("click", saveToFile, false);
refresh.addEventListener("click", convertFromText, false);
run.addEventListener("click", runDemo, false);

document.addEventListener("DOMContentLoaded", () => {
  //check too see if clide-build was started with any files
  getFiles();
});

function getFiles() {
  fetch("http://localhost:8080/getFiles").then(res => {
    if (res.status == 200) {
      res.json().then(json => {
        if (json.jsonText) {
          jsonOutput.innerText = JSON.stringify(
            JSON.parse(json.jsonText),
            null,
            4
          );
        }
        if (json.scriptText) {
          scriptOutput.innerText = json.scriptText;
        }
      });
    }
  });
}

function handleFiles(e) {
  const fileList = e.target.files;

  let reader = new FileReader();
  reader.onloadend = () => {
    jsonOutput.innerText = "File conversion in process";

    if (!fileList[0].name.includes(".json")) {
      scriptOutput.innerText = reader.result;
    }

    fetch("http://localhost:8080/convert", {
      method: "POST",
      header: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        filename: fileList[0].name,
        fileContents: reader.result
      })
    }).then(res => {
      if (res.status == 200) {
        res.json().then(json => {
          jsonOutput.innerText = JSON.stringify(json, null, 4);
        });
      } else {
        res.text().then(text => {
          jsonOutput.innerText = text;
        });
      }
    });
  };

  reader.readAsText(fileList[0]);
}

function saveToFile() {
  let textToSave = jsonOutput.value;
  let textToSaveAsBlob = new Blob([textToSave], { type: "text/plain" });
  let textToSaveAsURL = window.URL.createObjectURL(textToSaveAsBlob);
  let fileNameToSaveAs = "clide-demo.json";

  let downloadLink = document.createElement("a");
  downloadLink.download = fileNameToSaveAs;
  downloadLink.innerHTML = "Download File";
  downloadLink.href = textToSaveAsURL;
  downloadLink.onclick = e => {
    document.body.removeChild(e.target);
  };
  downloadLink.style.display = "none";
  document.body.appendChild(downloadLink);

  downloadLink.click();
}

function convertFromText() {
  fetch("http://localhost:8080/convert", {
    method: "POST",
    header: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      filename: "",
      fileContents: scriptOutput.value
    })
  }).then(res => {
    if (res.status == 200) {
      res.json().then(json => {
        jsonOutput.innerText = JSON.stringify(json, null, 4);
      });
    } else {
      res.text().then(text => {
        jsonOutput.innerText = text;
      });
    }
  });
}

function runDemo() {
  fetch("http://localhost:8080/run", {
    method: "POST",
    header: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      filename: "",
      fileContents: jsonOutput.value
    })
  });
}
