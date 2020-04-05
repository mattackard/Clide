//global config elements
const user = document.getElementById("user");
const directory = document.getElementById("directory");
const typeSpeed = document.getElementById("typespeed");
const humanize = document.getElementById("humanize");
const hideWarnings = document.getElementById("hideWarnings");
const clearBeforeAll = document.getElementById("clearBeforeAll");
const keyTriggerAll = document.getElementById("keyTriggerAll");
const fontPath = document.getElementById("fontPath");
const fontSize = document.getElementById("fontSize");
const windowContainer = document.getElementById("windowContainer");
const triggerKeys = document.getElementById("triggerKeys");
const keyList = document.getElementById("keyList");
const backgroundColor = document.getElementById("backgroundColor");
const primaryColor = document.getElementById("primaryColor");
const userColor = document.getElementById("userColor");
const directoryColor = document.getElementById("directoryColor");
const userPreview = document.getElementById("userPreview");
const directoryPreview = document.getElementById("directoryPreview");
const primaryText = document.getElementsByClassName("primaryText");
const addWindow = document.getElementById("addWindow");
const arrangeWindows = document.getElementById("arrangeWindows");
const recordKey = document.getElementById("recordKey");
let removeButtons = document.getElementsByClassName("removeButton");

//command config elements
const commands = document.getElementById("commands");
let addCommand = document.getElementsByClassName("addCommand");

const fileInput = document.getElementById("fileInput");
const saveFile = document.getElementById("saveFile");
const run = document.getElementById("run");

fileInput.addEventListener("change", handleFiles, false);
saveFile.addEventListener("click", saveToFile, false);
run.addEventListener("click", runDemo, false);

document.addEventListener("DOMContentLoaded", () => {
  //get any files passed in from cli or another page
  getFiles();
});

addWindow.addEventListener("click", () => {
  windowContainer.innerHTML += `<div class="clideWindow">
        <button class="removeButton" onclick="removeElement(this)">X</button>
        <div>
          <label for="windowName">Name</label>
          <input type="text" class="windowName" value="New Window" />
        </div>
        <div>
          <label for="x">X Position</label>
          <input type="number" class="x number" value="0" />
        </div>
        <div>
          <label for="y">Y Position</label>
          <input type="number" class="y number" value="0" />
        </div>
        <div>
          <label for="height">Vertical Resolution</label>
          <input type="number" class="height number" value="600" />
        </div>
        <div>
          <label for="width">Horizontal Resolution</label>
          <input type="number" class="width number" value="1000" /></div>`;
});

function addNewCommand() {
  commands.innerHTML += `<div class="command">
        <button class="removeButton" onclick="removeElement(this)">X</button>
        <input type="text" class="cmd" value="New Command" />
        <label for="window">Window</label><input type="text" class="window"/>
        <label for="predelay">PreDelay</label
        ><input type="number" class="predelay" placeholder="500" value="500" />
        <label for="postdelay">PostDelay</label
        ><input type="number" class="postdelay" placeholder="500" value="500" />
        <label for="timeout">Timeout</label
        ><input type="number" class="timeout" placeholder="500" />
        <div>
        <label for="typed">Typed</label><input type="checkbox" class="typed"/>
        <label for="hidden">Hidden</label><input type="checkbox" class="hidden"/>
        <label for="waitForKey">Wait for key press</label
        ><input type="checkbox" class="waitForKey"/>
        <label for="clearBeforeRun">Clear window before execution</label
        ><input type="checkbox" class="clearBeforeRun" />
        <label for="async">Asynchronous</label
        ><input type="checkbox" class="async"/></div>`;
}

function removeElement(element) {
  element.parentNode.outerHTML = "";
}

function getFiles() {
  fetch("http://localhost:8080/getFiles").then((res) => {
    if (res.status == 200) {
      res.json().then((json) => {
        if (json.jsonText) {
          populateConfig(JSON.parse(json.jsonText));
        }
      });
    }
  });
}

function handleFiles(e) {
  const fileList = e.target.files;

  let reader = new FileReader();
  reader.onloadend = () => {
    fetch("http://localhost:8080/convert", {
      method: "POST",
      header: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        filename: fileList[0].name,
        fileContents: reader.result,
      }),
    }).then((res) => {
      if (res.status == 200) {
        res.json().then((json) => {
          populateConfig(json);
        });
      } else {
        console.log("server error converting file to json");
      }
    });
  };

  reader.readAsText(fileList[0]);
}

function buildJSON() {
  //create an array containing all commands as json objects
  let commandArr = [];
  let commandList = document.getElementsByClassName("command");
  for (let i = 0; i < commandList.length; i++) {
    let cmd = commandList[i].getElementsByClassName("cmd").item(0);
    let typed = commandList[i].getElementsByClassName("typed").item(0);
    let window = commandList[i].getElementsByClassName("window").item(0);
    let predelay = commandList[i].getElementsByClassName("predelay").item(0);
    let postdelay = commandList[i].getElementsByClassName("postdelay").item(0);
    let timeout = commandList[i].getElementsByClassName("timeout").item(0);
    let hidden = commandList[i].getElementsByClassName("hidden").item(0);
    let waitForKey = commandList[i]
      .getElementsByClassName("waitForKey")
      .item(0);
    let clearBeforeRun = commandList[i]
      .getElementsByClassName("clearBeforeRun")
      .item(0);
    let async = commandList[i].getElementsByClassName("async").item(0);

    commandArr.push({
      cmd: cmd.value,
      typed: typed.checked,
      window: window.value,
      predelay: parseInt(predelay.value),
      postdelay: parseInt(postdelay.value),
      timeout: parseInt(timeout.value),
      hidden: hidden.checked,
      waitForKey: waitForKey.checked,
      clearBeforeRun: clearBeforeRun.checked,
      async: async.checked,
    });
  }

  //create an array containing window objects
  let windowArr = [];
  let windowList = document.getElementsByClassName("clideWindow");
  for (let i = 0; i < windowList.length; i++) {
    let name = windowList[i].getElementsByClassName("windowName").item(0);
    let x = windowList[i].getElementsByClassName("x").item(0);
    let y = windowList[i].getElementsByClassName("y").item(0);
    let height = windowList[i].getElementsByClassName("height").item(0);
    let width = windowList[i].getElementsByClassName("width").item(0);

    windowArr.push({
      name: name.value,
      x: parseInt(x.value),
      y: parseInt(y.value),
      height: parseInt(height.value) > 0 ? parseInt(height.value) : 600,
      width: parseInt(width.value) > 0 ? parseInt(width.value) : 1000,
    });
  }

  //create an array of key strings
  let keyArr = keyList.innerText.split(" ");

  //createa new clide demo json and populate with all form fields
  let newClide = {
    user: user.value,
    directory: directory.value,
    typespeed: parseInt(typeSpeed.value),
    humanize: parseFloat(humanize.value),
    hideWarnings: hideWarnings.checked,
    clearBeforeAll: clearBeforeAll.checked,
    keyTriggerAll: keyTriggerAll.checked,
    fontPath: fontPath.value,
    fontSize: parseInt(fontSize.value),
    colorScheme: {
      userText: "0,150,255,255",
      directoryText: "150,255,150,255",
      primaryText: "220,220,220,255",
      terminalBG: "30,30,30,255",
    },
    windows: windowArr,
    triggerKeys: keyArr,
    commands: commandArr,
  };
  return newClide;
}

function saveToFile() {
  let textToSave = buildJSON();
  let textToSaveAsBlob = new Blob([textToSave], { type: "text/plain" });
  let textToSaveAsURL = window.URL.createObjectURL(textToSaveAsBlob);
  let fileNameToSaveAs = "clide-demo.json";

  let downloadLink = document.createElement("a");
  downloadLink.download = fileNameToSaveAs;
  downloadLink.innerHTML = "Download File";
  downloadLink.href = textToSaveAsURL;
  downloadLink.onclick = (e) => {
    document.body.removeChild(e.target);
  };
  downloadLink.style.display = "none";
  document.body.appendChild(downloadLink);

  downloadLink.click();
}

function populateConfig(clide) {
  user.value = clide.user;
  directory.value = clide.directory;
  typeSpeed.value = clide.typespeed;
  humanize.value = clide.humanize;
  hideWarnings.checked = clide.hideWarnings;
  clearBeforeAll.checked = clide.clearBeforeAll;
  keyTriggerAll.checked = clide.keyTriggerAll;
  fontPath.value = clide.fontPath;
  fontSize.value = clide.fontSize;

  //apply color values to inputs with hex conversion
  if (clide.colorScheme) {
    backgroundColor.value = byteToHex(clide.colorScheme.terminalBG);
    primaryColor.value = byteToHex(clide.colorScheme.primaryText);
    userColor.value = byteToHex(clide.colorScheme.userText);
    directoryColor.value = byteToHex(clide.colorScheme.directoryText);
  }

  //build window html
  if (clide.windows) {
    let html = "";
    clide.windows.forEach((window) => {
      html += `<div class="clideWindow">
        <button class="removeButton" onclick="removeElement(this)">X</button>
        <div>
          <label for="windowName">Name</label>
          <input type="text" class="windowName" value="${window.name}" />
        </div>
        <div>
          <label for="x">X Position</label>
          <input type="number" class="x number" value="${window.x}" />
        </div>
        <div>
          <label for="y">Y Position</label>
          <input type="number" class="y number" value="${window.y}" />
        </div>
        <div>
          <label for="height">Vertical Resolution</label>
          <input type="number" class="height number" value="${window.height}" />
        </div>
        <div>
          <label for="width">Horizontal Resolution</label>
          <input type="number" class="width number" value="${window.width}" /></div></div >`;
    });
    windowContainer.innerHTML = html;
  }

  //add trigger keys to span
  if (clide.triggerKeys) {
    keyText = "";
    clide.triggerKeys.forEach((key) => {
      keyText += key + " ";
    });
    keyList.innerText = keyText;
  }

  cmdHTML = `<h1>Command Configuration</h1><button class="addCommand" onclick="addNewCommand()">Add Commmand</button>`;
  clide.commands.forEach((command) => {
    cmdHTML += `<div class="command">
        <button class="removeButton" onclick="removeElement(this)">X</button>
        <input type="text" class="cmd" value="${command.cmd}" />
        <label for="window">Window</label><input type="text" class="window" value="${
          command.window
        }" />
        <label for="predelay">PreDelay</label
        ><input type="number" class="predelay" placeholder="500" value="${
          command.predelay
        }" />
        <label for="postdelay">PostDelay</label
        ><input type="number" class="postdelay" placeholder="500" value="${
          command.postdelay
        }" />
        <label for="timeout">Timeout</label
        ><input type="number" class="timeout" placeholder="500" value="${
          command.timeout
        }" />
        <div>
        <label for="typed">Typed</label><input type="checkbox" class="typed" ${
          command.typed ? "checked" : ""
        } />
        <label for="hidden">Hidden</label><input type="checkbox" class="hidden" ${
          command.hidden ? "checked" : ""
        } />
        <label for="waitForKey">Wait for key press</label
        ><input type="checkbox" class="waitForKey" ${
          command.waitForKey ? "checked" : ""
        } />
        <label for="clearBeforeRun">Clear window before execution</label
        ><input type="checkbox" class="clearBeforeRun" ${
          command.clearBeforeRun ? "checked" : ""
        } />
        <label for="async">Asynchronous</label
        ><input type="checkbox" class="async" ${
          command.async ? "checked" : ""
        } /></div>
      </div>`;
  });

  commands.innerHTML = cmdHTML;
}

//takes an rgba in the form of 255,255,255,255 and return a hex value
function byteToHex(byteString) {
  split = byteString.split(",");
  while (split.length > 3) {
    split.pop();
  }
  hex = "#";
  split.forEach((byte) => {
    num = parseInt(byte);
    if (byte == 0) {
      hex += "0" + parseInt(byte).toString(16);
    } else {
      hex += parseInt(byte).toString(16);
    }
  });
  return hex;
}

//takes a hex value and returns an rgba in the form of 255,255,255,255
function hexToByte(hexString) {
  bytes = [];
  split = [
    hexString.substring(1, 3),
    hexString.substring(3, 5),
    hexString.substring(5, 7),
  ];
  split.forEach((hex) => {
    bytes.push(parseInt(hex, 16));
  });
  return bytes.join(",") + ",255";
}

function runDemo() {
  fetch("http://localhost:8080/run", {
    method: "POST",
    header: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      filename: "",
      fileContents: JSON.stringify(buildJSON(), null, 4),
    }),
  });
}
