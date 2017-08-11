// Подстановка имени группы в поле в модальном окне
// Для кнопки "Переименовать"
$('#editGroup').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);         // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');     // Извлечь информацию из "data-*" у кнопки
    var modal = $(this);                        // Обновить модальное окно
    modal.find('#id_old_group').val(recipientName);      // В input со старым значением
});

// Для кнопки "Удалить"
$('#delGroup').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);
    var recipientName = button.data('name');
    var modal = $(this);
    modal.find('.modal-body input').val(recipientName)      // Заполнить все input-ы (он один)
});

// Для кнопки "Добавить"
$('#addSuite').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);
    var recipientName = button.data('name');
    var modal = $(this);
    modal.find('#id_suites_group').val(recipientName)
});
