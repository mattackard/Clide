// global config elements
const user = document.getElementById("user");
const directory = document.getElementById("directory");
const typeSpeed = document.getElementById("typespeed");
const humanize = document.getElementById("humanize");
const hideWarnings = document.getElementById("hideWarnings");
const hideWindows = document.getElementById("hideWindows");
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
const colorArr = [backgroundColor, primaryColor, userColor, directoryColor];
const colorPreview = document.getElementById("colorPreview");
const userPreview = document.getElementById("userPreview");
const directoryPreview = document.getElementById("directoryPreview");
const primaryText = document.getElementsByClassName("primaryText");
const addWindow = document.getElementById("addWindow");
const recordKey = document.getElementById("recordKey");
let removeButtons = document.getElementsByClassName("removeButton");

// command config elements
const commands = document.getElementById("commands");
let addCommand = document.getElementsByClassName("addCommand");
let resizeWindows = document.getElementsByClassName("resizeWindows");

const fileInput = document.getElementById("fileInput");
const saveFile = document.getElementById("saveFile");
const run = document.getElementById("run");

fileInput.addEventListener("change", handleFiles, false);
saveFile.addEventListener("click", saveToFile, false);
run.addEventListener("click", runDemo, false);

document.addEventListener("DOMContentLoaded", () => {
  // get any files passed in from cli or another page
  getFiles();
});

// saves the json contents every time an element loses focus
document.addEventListener("focusout", () => {
  fetch("http://localhost:8080/save", {
    method: "POST",
    header: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      jsonText: JSON.stringify(buildJSON()),
    }),
  }).then((res) => {
    if (res.status != 200) {
      console.log("server error converting file to json");
    }
  });
});

addWindow.addEventListener("click", () => {
  let newHTML = `<div class="clideWindow">
  <button class="removeButton" onclick="removeElement(this.parentNode)">X</button>
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
  windowContainer.insertAdjacentHTML("beforeend", newHTML);
});

recordKey.addEventListener("click", () => {
  recordKey.innerText = "Press a key";
  document.addEventListener("keydown", listenForOneKey);
});

colorArr.forEach((color) => {
  color.addEventListener("change", (e) => {
    switch (e.target.id) {
      case "backgroundColor":
        colorPreview.style.backgroundColor = e.target.value;
        break;
      case "primaryColor":
        for (let text of primaryText) {
          text.style.color = e.target.value;
        }
        break;
      case "userColor":
        userPreview.style.color = e.target.value;
        break;
      case "directoryColor":
        directoryPreview.style.color = e.target.value;
        break;
    }
  });
});

function addNewCommand() {
  commands.insertAdjacentHTML(
    "beforeend",
    `<div class="command">
        <button class="removeButton" onclick="removeElement(this.parentNode)">X</button>
        <button class="moveUp" onclick="swapCommand(this)">&#8593;</button>
        <button class="moveDown" onclick="swapCommand(this)">&#8595;</button>
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
          ><input type="checkbox" class="async"/>
          <label for="resizeWindows">Resize windows</label
          ><input type="checkbox" class="resizeWindows" onchange="showWindows(this)"/>
        </div>
        <div class="resize"></div>
    </div>`
  );
}

function arrangeWindows(element) {
  fetch("http://localhost:8080/arrangeWindows", {
    method: "POST",
    header: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      fileContents: JSON.stringify(buildWindowJSON(element.parentNode)),
    }),
  }).then(res => res.json()).then(json => {
    let newHTML = "";
    json.windows.forEach(window => {
      newHTML += `<div class="clideWindow">
      <button class="removeButton" onclick="removeElement(this.parentNode)">X</button>
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
        <input type="number" class="width number" value="${window.width}" /></div></div >`
    });
    element.previousElementSibling.innerHTML = newHTML;
  });
}

function removeElement(element) {
  element.outerHTML = null;
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
  // create an array containing all commands as json objects
  let commandArr = [];
  let commandList = document.getElementsByClassName("command");
  for (let i = 0; i < commandList.length; i++) {
    let command = storeCommandValues(commandList[i]);
    commandArr.push(command);
  }

  // create an array containing window objects
  let windowArr = buildWindowJSON(windowContainer).windows;

  // create an array of key strings
  let keyArr = [];
  keyList.childNodes.forEach((key) => {
    keyArr.push(key.innerText.substring(0, key.innerText.length - 1));
  });

  // create a new clide demo json and populate with all form fields
  let newClide = {
    user: user.value,
    directory: directory.value,
    typespeed: parseInt(typeSpeed.value),
    humanize: parseFloat(humanize.value),
    hideWarnings: hideWarnings.checked,
    hideWindows: hideWindows.checked,
    clearBeforeAll: clearBeforeAll.checked,
    keyTriggerAll: keyTriggerAll.checked,
    fontPath: fontPath.value,
    fontSize: parseInt(fontSize.value),
    colorScheme: {
      userText: hexToByte(userColor.value),
      directoryText: hexToByte(directoryColor.value),
      primaryText: hexToByte(primaryColor.value),
      terminalBG: hexToByte(backgroundColor.value),
    },
    windows: windowArr,
    triggerKeys: keyArr,
    commands: commandArr,
  };
  return newClide;
}

function storeCommandValues(command) {
  let cmd = command.getElementsByClassName("cmd").item(0);
  let typed = command.getElementsByClassName("typed").item(0);
  let window = command.getElementsByClassName("window").item(0);
  let predelay = command.getElementsByClassName("predelay").item(0);
  let postdelay = command.getElementsByClassName("postdelay").item(0);
  let timeout = command.getElementsByClassName("timeout").item(0);
  let hidden = command.getElementsByClassName("hidden").item(0);
  let waitForKey = command
    .getElementsByClassName("waitForKey")
    .item(0);
  let clearBeforeRun = command
    .getElementsByClassName("clearBeforeRun")
    .item(0);
  let async = command.getElementsByClassName("async").item(0);
  let resizedWindows = buildWindowJSON(command).windows;

  return {
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
    resizeWindows: resizedWindows
  };
}

function updateCommandValues(command, newValues) {
  let cmd = command.getElementsByClassName("cmd").item(0);
  let typed = command.getElementsByClassName("typed").item(0);
  let window = command.getElementsByClassName("window").item(0);
  let predelay = command.getElementsByClassName("predelay").item(0);
  let postdelay = command.getElementsByClassName("postdelay").item(0);
  let timeout = command.getElementsByClassName("timeout").item(0);
  let hidden = command.getElementsByClassName("hidden").item(0);
  let waitForKey = command
    .getElementsByClassName("waitForKey")
    .item(0);
  let clearBeforeRun = command
    .getElementsByClassName("clearBeforeRun")
    .item(0);
  let async = command.getElementsByClassName("async").item(0);
  let resizeWindows = command.getElementsByClassName("resizeWindows").item(0);
  let resizeDiv = command.getElementsByClassName("resize").item(0);
  let arrangeBtn = command.getElementsByClassName("arrange");


    cmd.value = newValues.cmd;
    typed.checked = newValues.typed;
    window.value = newValues.window;
    predelay.value = newValues.predelay;
    postdelay.value = newValues.postdelay;
    timeout.value = newValues.timeout;
    hidden.checked = newValues.hidden;
    waitForKey.checked = newValues.waitForKey;
    clearBeforeRun.checked = newValues.clearBeforeRun;
    async.checked = newValues.checked;
    resizeWindows.checked = newValues.resizeWindows.length > 0;

    //remove any arrange windows buttons that might be left over
    if (arrangeBtn.length > 0) {
      for (let btn of arrangeBtn) {
        removeElement(btn);
      }
    }

    //populate the resize div with windows for resizing
    let resizeHTML = ``;
    if (resizeWindows.checked) {
      newValues.resizeWindows.forEach(window => {
        resizeHTML += `<div class="clideWindow">
          <button class="removeButton" onclick="removeElement(this.parentNode)">X</button>
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
            <input type="number" class="width number" value="${window.width}" /></div></div>`;
      });  
      resizeDiv.innerHTML = resizeHTML;
      resizeDiv.outerHTML += `<button class="arrange" onclick="arrangeWindows(this)">Arrange Windows</button>`;
    } else {
      resizeDiv.innerHTML = null;
    }
}

function buildWindowJSON(container) {
  // create an array containing window objects
  let windowArr = [];
  let windowList = container.getElementsByClassName("clideWindow");
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

  return {windows: windowArr};
}

// saveToFile saves the json file using the browsers default behavior
function saveToFile() {
  let textToSave = JSON.stringify(buildJSON(), null, 4);
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

// populateConfig populates all the input fields in the editor window
function populateConfig(clide) {
  user.value = clide.user;
  directory.value = clide.directory;
  typeSpeed.value = clide.typespeed;
  humanize.value = clide.humanize;
  hideWarnings.checked = clide.hideWarnings;
  hideWindows.checked = clide.hideWindows;
  clearBeforeAll.checked = clide.clearBeforeAll;
  keyTriggerAll.checked = clide.keyTriggerAll;
  fontPath.value = clide.fontPath ? clide.fontPath : "";
  fontSize.value = clide.fontSize;

  // apply color values to inputs with hex conversion
  // set to black and white if no colorscheme is set
  if (clide.colorScheme) {
    let bg = clide.colorScheme.terminalBG
      ? byteToHex(clide.colorScheme.terminalBG)
      : "#000000";
    backgroundColor.value = bg;
    colorPreview.style.backgroundColor = bg;

    let pc = clide.colorScheme.primaryText
      ? byteToHex(clide.colorScheme.primaryText)
      : "#FFFFFF";
    primaryColor.value = pc;
    for (let text of primaryText) {
      text.style.color = pc;
    }

    let uc = clide.colorScheme.userText
      ? byteToHex(clide.colorScheme.userText)
      : "#FFFFFF";
    userColor.value = uc;
    userPreview.style.color = uc;

    let dc = clide.colorScheme.directoryText
      ? byteToHex(clide.colorScheme.directoryText)
      : "#FFFFFF";
    directoryColor.value = dc;
    directoryPreview.style.color = dc;
  } else {
    backgroundColor.value = "#000000";
    primaryColor.value = "#FFFFFF";
    userColor.value = "#FFFFFF";
    directoryColor.value = "#FFFFFF";
  }

  // build window html
  if (clide.windows) {
    let html = "";
    clide.windows.forEach((window) => {
      html += `<div class="clideWindow">
        <button class="removeButton" onclick="removeElement(this.parentNode)">X</button>
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
    windowHTML = html
    windowContainer.innerHTML = html;
  }

  // add trigger keys to list
  if (clide.triggerKeys) {
    let keyText = "";
    clide.triggerKeys.forEach((key) => {
      keyText += `<li class="triggerKey">${key}<button class="removeButtonSmall" onclick="removeElement(this.parentNode)">
              X
            </button></li>`;
    });
    keyList.innerHTML = keyText;
  }

  // build all command divs
  cmdHTML = `<h1>Command Configuration</h1><button class="addCommand" onclick="addNewCommand()">Add Commmand</button>`;
  clide.commands.forEach(command => {
    let resizeHTML = "";
    if (command.resizeWindows) {
      command.resizeWindows.forEach(window => {
        resizeHTML += `<div class="clideWindow">
          <button class="removeButton" onclick="removeElement(this.parentNode)">X</button>
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
            <input type="number" class="width number" value="${window.width}" /></div></div >`
      });
    }
    

    cmdHTML += `<div class="command">
        <button class="removeButton" onclick="removeElement(this.parentNode)">X</button>
        <button class="moveUp" onclick="swapCommand(this)">&#8593;</button>
        <button class="moveDown" onclick="swapCommand(this)">&#8595;</button>
        <input type="text" class="cmd" value="${command.cmd.replace(/"/g, '&quot;')}" />
        <label for="window">Window</label><input type="text" class="window" value="${
          command.window ? command.window : ""
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
        ><input type="number" class="timeout" placeholder="5" value="${
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
        } />
        <label for="resizeWindows">Resize windows</label
        ><input type="checkbox" class="resizeWindows" onchange="showWindows(this)" ${
          resizeHTML == "" ? "" : "checked"
        } /></div>
        <div class="resize">${resizeHTML == "" ? "" : resizeHTML}</div>
        ${resizeHTML == "" ? "" : `<button class="arrange" onclick="arrangeWindows(this)">Arrange Windows</button>`}
      </div>`;
  });

  commands.innerHTML = cmdHTML;
}

// show windows creates a div with all windows for resizing in a command
function showWindows(element) {
  let div = element.parentNode.parentNode.getElementsByClassName("resize");
  if (element.checked) {
    div[0].innerHTML = windowContainer.innerHTML;
    div[0].outerHTML += `<button class="arrange" onclick="arrangeWindows(this)">Arrange Windows</button>`;
  } else {
    div[0].innerHTML = null;
    div[0].nextSibling.outerHTML = null;
  }
}

// takes an rgba in the form of 255,255,255,255 and return a hex value
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

// takes a hex value and returns an rgba in the form of 255,255,255,255
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

// sends the json data to run with clide
function runDemo() {
  fetch("http:// localhost:8080/run", {
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

// listens for a single key press and then removes the event listener
function listenForOneKey(event) {
  document.removeEventListener(event.type, listenForOneKey);
  let key = event.key;
  if (key.startsWith("Arrow")) {
    key = key.substring(5);
  } else if (key == " ") {
    key = "Space";
  }

  if (key.length == 1) {
    key = key.toUpperCase();
  }
  keyList.innerHTML += `<li class="triggerKey">${key}<button class="removeButtonSmall" onclick="removeElement(this.parentNode)">
        X
      </button></li>`;
  recordKey.innerText = "Record Keypress";
}

// swapCommand swaps two commands depending on if the move up or move down buttons are pressed
function swapCommand(element) {
  //store the input values of the passed in command
  let thisCommand = storeCommandValues(element.parentNode);

  if (element.classList.contains("moveUp")) {
    if (!element.parentNode.previousElementSibling.classList.contains("addCommand")) {
      //store the input values for the previous command
      let previousCommand = storeCommandValues(element.parentNode.previousElementSibling);

      //swap commands
      updateCommandValues(element.parentNode, previousCommand);
      updateCommandValues(element.parentNode.previousElementSibling, thisCommand);
    }
  } else if (element.parentNode.nextSibling) {
    //store the input values for the next command
    let nextCommand = storeCommandValues(element.parentNode.nextSibling);

    //swap commands
    updateCommandValues(element.parentNode, nextCommand);
    updateCommandValues(element.parentNode.nextSibling, thisCommand);
  }
}
