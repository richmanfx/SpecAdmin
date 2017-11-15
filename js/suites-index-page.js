
// Скрипт для выпадающих списков
function anichange (objName) {
    if ( $(objName).css('display') === 'none' ) {
        $(objName).animate({height: 'show'}, 400);
    } else {
        $(objName).animate({height: 'hide'}, 200);
    }
}
function anichange(a) {
    $(a).slideToggle(400);
}


/// Подстановка имени сюиты в поле в модальном окне
// Для кнопки "Переименовать Сюиту"
$('#renameSuite').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);         // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');     // Извлечь информацию из "data-name" у кнопки
    var modal = $(this);                         // Обновить модальное окно
    modal.find('#id_old_suite').val(recipientName)      // Только в input со старым значением
});

// Для кнопки "Добавить Сюиту"
$('#addSuite').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);
    var recipientName = button.data('name');
    var modal = $(this);
    modal.find('#id_suites_group').val(recipientName)
});

// Для кнопки "Редактировать Сюиту"
$('#editSuite').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);
    var recipientName = button.data('name');
    var modal = $(this);
    modal.find('#id_suite').val(recipientName)
});

// Для кнопки "Удалить Сюиту"
$('#delSuite').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);
    var recipientName = button.data('name');
    var modal = $(this);
    modal.find('#id_suite').val(recipientName)
});

// Для кнопки "Добавить Сценарий"
$('#addScript').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);
    var recipientName = button.data('name');
    var modal = $(this);
    modal.find('#id_scripts_suite').val(recipientName)
});


// Подстановка имени Сценария и Сюиты в поля в модальном окне
// Для кнопки "Редактировать"
$('#editScript').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);         // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');     // Извлечь информацию из "data-name" у кнопки
    var recipientSuite = button.data('suite');     // Извлечь информацию из "data-suite" у кнопки
    var modal = $(this);         // Обновить модальное окно
    modal.find('#id_script').val(recipientName);      // В input с именем сценария
    modal.find('#id_scripts_suite').val(recipientSuite)      // В input с именем сюиты
});

// Для кнопки "Удалить"
$('#delScript').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);         // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');     // Извлечь информацию из "data-name" у кнопки
    var recipientSuite = button.data('suite');     // Извлечь информацию из "data-suite" у кнопки
    var modal = $(this);         // Обновить модальное окно
    modal.find('#id_script').val(recipientName);
    modal.find('#id_scripts_suite').val(recipientSuite)      // В input с именем сюиты
});

// Для кнопки "Добавить Шаг"
$('#addStep').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);
    var recipientName = button.data('name');
    var recipientSuite = button.data('suite');
    var modal = $(this);
    modal.find('#id_steps_script').val(recipientName);
    modal.find('#id_scripts_suite').val(recipientSuite)
});


// Подстановка имени Шага, его Сценария и Сюиты в поля в модальном окне
// Для кнопки "Редактировать" Шага
$('#editStep').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);            // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');        // Извлечь информацию из "data-name" у кнопки
    var recipientScriptsId = button.data('script');     // Извлечь "ScriptsId"

    var stepsScriptName = "";
    var scripsSuiteName = "";

    // Получить имя Сценария по его ScriptsId и имя Сюиты, которой принадлежит Сценарий
    $.ajax({
        async: false,
        type: 'POST',
        url: '/spec-admin/get-steps-options',
        data: 'ScriptsId=' + recipientScriptsId,
        success: function(answerFromServer){
            // alert(answerFromServer.stepsScriptName + " и " + answerFromServer.scripsSuiteName);
            stepsScriptName = answerFromServer.stepsScriptName;
            scripsSuiteName = answerFromServer.scripsSuiteName;
        },
        error: function(){
            alert("Ошибка при ответе на AJAX POST запрос имени Сценария по его ScriptsId и имени Сюиты.");
        }

    });

    var modal = $(this);         // Обновить модальное окно
    modal.find('#id_step').val(recipientName);                  // В input с именем Шага
    modal.find('#id_steps_script').val(stepsScriptName);        // В input с именем Сценария
    modal.find('#id_scripts_suite').val(scripsSuiteName);       // В input с именем Сюиты

});

// Для кнопки "Копировать" Шага
$('#copyStep').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);
    // var name = button.data('name');
    var stepId = button.data('stepid');

    // Поместить шаг в буфер копирования
    $.ajax({
        async: false,
        type: 'POST',
        url: '/spec-admin/copy-step-in-clipboard',
        data: 'StepId=' + stepId,
        success: function(answerFromServer){
            // alert(answerFromServer);
        },
        error: function(){
            alert("Ошибка при ответе на AJAX POST запрос копирования Шага в буфер обмена");
        }

    });

});

// Для кнопки "Вставить" Шага
$('#insertStep').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);                // Кнопка, вызвавшая модальное окно
    var recipientsStepId = button.data('stepid');       // Извлечь информацию из "data-stepid" у кнопки
    var recipientScriptsId = button.data('script');     // Извлечь из "data-script" у кнопки

    // Получить параметры нового шага из буфера обмена - имя, описание и ожидаемый результат шага
    $.ajax({
        async: false,
        type: 'POST',
        url: '/spec-admin/get-step-from-buffer',
        data: 'StepsId=' + recipientsStepId,
        success: function(answerFromServer){
            // alert(
            //     answerFromServer.stepsName + " и " +
            //     answerFromServer.stepsDescription + " и " +
            //     answerFromServer.stepsExpectedResult
            // );
            stepsName = answerFromServer.stepsName;
            stepsDescription = answerFromServer.stepsDescription;
            stepsExpectedResult = answerFromServer.stepsExpectedResult;
        },
        error: function(answerFromServer){
            alert("Ошибка при ответе на AJAX POST запрос параметров Шага из буфера обмена: " + answerFromServer.Status);
        }
    });


    // Получить имя Сценария по его ScriptsId и имя Сюиты, которой принадлежит Сценарий
    $.ajax({
        async: false,
        type: 'POST',
        url: '/spec-admin/get-steps-options',
        data: 'ScriptsId=' + recipientScriptsId,
        success: function(answerFromServer){
            // alert(answerFromServer.stepsScriptName + " и " + answerFromServer.scripsSuiteName);
            stepsScriptName = answerFromServer.stepsScriptName;
            scripsSuiteName = answerFromServer.scripsSuiteName;
        },
        error: function(){
            alert("Ошибка при ответе на AJAX POST запрос имени Сценария и имени Сюиты Шага по ScriptsId Шага.");
        }

    });

    var modal = $(this);         // Обновить модальное окно
    // modal.find('#id_step').val(recipientName);                  // В input с именем Шага
    modal.find('#id_steps_script').val(stepsScriptName);        // В input с именем Сценария
    modal.find('#id_scripts_suite').val(scripsSuiteName);       // В input с именем Сюиты
    modal.find('#id_step').val(stepsName);      // В textarea c именем Шага
    modal.find('#id_steps_description').val(stepsDescription)       // В textarea c описанием Шага
    modal.find('#id_steps_expected_result').val(stepsExpectedResult)    // В textarea с ожидаемым результатом Шага
});

// Для кнопки "Удалить" Шага
$('#delStep').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);         // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');     // Извлечь информацию из "data-name" у кнопки
    var recipientScriptsId = button.data('script');     // Извлечь "ScriptsId"

    var stepsScriptName = "";
    var scripsSuiteName = "";

    // Получить имя Сценария по его ScriptsId и имя Сюиты, которой принадлежит Сценарий
    $.ajax({
        async: false,
        type: 'POST',
        url: '/spec-admin/get-steps-options',
        data: 'ScriptsId=' + recipientScriptsId,
        success: function(answerFromServer){
            // alert(answerFromServer.stepsScriptName + " и " + answerFromServer.scripsSuiteName);
            stepsScriptName = answerFromServer.stepsScriptName;
            scripsSuiteName = answerFromServer.scripsSuiteName;
        },
        error: function(){
            alert("Ошибка при ответе на AJAX POST запрос удаления Шага");
        }

    });

    var modal = $(this);         // Обновить модальное окно
    modal.find('#id_step').val(recipientName);
    modal.find('#id_script').val(stepsScriptName);      // В input с именем Сценария
    modal.find('#id_suite').val(scripsSuiteName);       // В input с именем Сюиты
});


/// Печать на принтер

// Печатать Шаги Сценария
$('#printSteps').on('show.bs.modal', function (event) {
        var button = $(event.relatedTarget);        // Кнопка, вызвавшая модальное окно
        var scriptName = button.attr("data-name");    // Извлечь информацию из "data-name" у кнопки
        var suiteName = button.attr("data-suite");    // Извлечь информацию из "data-suite" у кнопки

        $.ajax({
            async: false,
            type: 'POST',
            url: '/spec-admin/print-scripts-steps',
            data: 'scriptName=' + scriptName + ';' + 'suiteName=' + suiteName,
            success: function(answerFromServer){
                console.log("success: " + answerFromServer);
                // alert(answerFromServer);
            },
            error: function(){
                // alert("Ошибка при ответе на AJAX POST запрос для печати ПДФ");
            }
        });
});

// Печатать Сценарии Сюиты
$('#printScripts').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);        // Кнопка, вызвавшая модальное окно
    var suiteName = button.attr("data-name");    // Извлечь информацию из "data-name" у кнопки

    $.ajax({
        async: false,
        type: 'POST',
        url: '/spec-admin/print-suites-scripts',
        data: 'suiteName=' + suiteName,
        success: function(answerFromServer){
            console.log("success: " + answerFromServer);
            // alert(answerFromServer);
        },
        error: function(){
            // alert("Ошибка при ответе на AJAX POST запрос для печати ПДФ");
        }
    });
});














