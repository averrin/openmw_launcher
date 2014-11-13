import QtQuick 2.2
import QtQuick.Controls 1.1
import QtQuick.Layouts 1.0

ApplicationWindow {
    visible: true
    title: "OpenMW Launcher"
    property int margin: 11
    width: mainLayout.implicitWidth + 2 * margin
    height: mainLayout.implicitHeight + 2 * margin
    minimumWidth: mainLayout.Layout.minimumWidth + 2 * margin + 40
    minimumHeight: mainLayout.Layout.minimumHeight + 2 * margin
    
    RowLayout {
            id: mainLayout
            anchors.fill: parent
            anchors.margins: margin

        ColumnLayout {
            anchors.fill: parent
            anchors.margins: margin
    
            Label {text: "Installed version: " + localVersion}
            Label {
                id: rlabel
                objectName: "Rlabel"
                text: "Fetching new version..." 
            }
            
            
            ComboBox {
                id: profilesBox
               property var pseudoModel: []
               model: pseudoModel
               Component.onCompleted: {
                   var i;
                   for(i=0; i < profiles.length(); i++){
                       pseudoModel.push(profiles.at(i))
                   }
                   model = pseudoModel
                   currentIndex = CurrentProfile
                   console.log(CurrentProfile)
               }
               onCurrentIndexChanged: {
                    profiles.select(currentIndex)
                    console.log("Changed")
               }
            }
            
            
            Button {
                text: "Play"
                onClicked: startOpenMW()
            }
    
    
        }
        Rectangle {
            width: 180; height: 200
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            anchors.margins: margin
            anchors.right: parent.right
            
            ListView {
                width: 120;
                model: contentFiles.length
                delegate: Text {
                    text: contentFiles.text(index)
                }
                anchors.top: parent.top
                anchors.bottom: parent.bottom
                anchors.margins: margin
                
            }
        }
    }
}
