import QtQuick 2.2
import QtQuick.Controls 1.1
import QtQuick.Layouts 1.0

ApplicationWindow {
    visible: true
    title: "OpenMW Launcher"
    property int margin: 11
    width: mainLayout.implicitWidth + 2 * margin
    height: mainLayout.implicitHeight + 2 * margin
    minimumWidth: mainLayout.Layout.minimumWidth + 2 * margin
    minimumHeight: mainLayout.Layout.minimumHeight + 2 * margin
    objectName: "Rlabel"

    ColumnLayout {
        id: mainLayout
        anchors.fill: parent
        anchors.margins: margin

        Label {text: "Installed version: " + localVersion}
        Label {
            id: rlabel
            text: "Available version: " + remoteVersion
        }
        
        
        ComboBox {
           property var pseudoModel: []
           model: pseudoModel
           Component.onCompleted: {
               var i;
               for(i=0; i<ProfilesModel.length(); i++){
                   pseudoModel.push(ProfilesModel.at(i))
               }
               model = pseudoModel
               currentIndex = CurrentProfile
           }
           onCurrentIndexChanged: ProfilesModel.select(currentIndex)
        }
        
        Button {
            text: "Launch"
            onClicked: startOpenMW()
        }


    }
}
