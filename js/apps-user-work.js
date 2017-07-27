/**
 * Created by Александр Ящук (R5AM, Zoer) on 27.07.2017.
 */


/// Подстановка Логина и Полного имени пользователя в поля в модальном окне

// Для кнопки "Удалить"
$('#deleteUser').on('show.bs.modal', function (event) {

    // Кнопка, вызвавшая модальное окно
    var button = $(event.relatedTarget);

    // Извлечь информацию из "data-*" полей у кнопки
    var deletedLogin = button.data('login');
    var deletedFullName = button.data('name');

    // Обновить модальное окно
    var modal = $(this);

    // В input-ы вставить
    modal.find('#id_login').val(deletedLogin)
    modal.find('#id_full_name').val(deletedFullName)
});

